package parser

import (
	"fmt"
	"github.com/huxiaoyugo/huxykit/hipl/ast"
	"github.com/huxiaoyugo/huxykit/hipl/err"
	"github.com/huxiaoyugo/huxykit/hipl/scanner"
	"github.com/huxiaoyugo/huxykit/hipl/token"
	"io/ioutil"
	"strconv"
	"strings"
)

type Parser struct {
	protoPath []string
	scan      []*scanner.Scanner
	level     int
	tokens    []*scanner.TokenInfo
	curToken  *scanner.TokenInfo
	nextToken *scanner.TokenInfo

	imports map[string]*Import

	// 是否需要检查字段命名是否合法
	needCheckIdentifier bool

	requestPool  map[string]*ast.Request
	resourcePool map[string]*ast.Resource
	structPools  map[string]*ast.Struct
	enumPools    map[string]*ast.Enum

	// 需要最后再确认是否存在的struct或者enum
	//structNeedCheckList []*scanner.TokenInfo
	//enumNeedCheckList   [] *scanner.TokenInfo
	// tmp
	typ ast.Type
}

func (p *Parser) Start() (e error) {
	defer func() {
		if er := recover(); er != nil {
			switch er.(type) {
			case error:
				e = er.(error)
			default:
				e = fmt.Errorf("%v", er)
			}
		}
	}()
	p.tokenPointInit()
	for p.nextToken.Tok != token.EOF {
		switch p.nextToken.Tok {
		case token.STR:
			switch getStringVal(p.nextToken.Lit) {
			case "request":
				p.parseRequest()
			case "resource":
				p.parseResource()
			default:
				p.errorInfo("invalid char")
			}
		case token.STRUCT:
			p.parseStruct()
		case token.ENUM:
			p.parseEnum()
		case token.PACKAGE:
			p.parsePackage()
		case token.IMPORT:
			p.parseImport()
		default:
			p.errorInfo("invalid char")
		}
	}
	return nil
}

func (p *Parser) isIdentifier(ident string) {
	if p.needCheckIdentifier {
		if !token.IsIdentifier(ident) {
			p.errorInfo(fmt.Sprintf("%s is invalid identifier", ident))
		}
	}
}

func (p *Parser) parseStruct() {
	stru := &ast.Struct{
		Tok: p.nextToken,
	}
	p.NextTokenMust(token.STRUCT)
	tokInfo := p.NextTokenMust(token.STR)
	stru.Name = getStringVal(tokInfo.Lit)
	p.isIdentifier(stru.Name)
	p.NextTokenMust(token.COLON)
	stru.Fields = p.parseStructBody()

	p.AddStructToPool(stru)
}

func (p *Parser) parseStructBody() []*ast.StructField {
	var fields []*ast.StructField
	p.NextTokenMust(token.LBRACE)
	for p.nextToken.Tok != token.RBRACE {
		fields = append(fields, p.parseStructField())
	}
	p.NextTokenMust(token.RBRACE)
	return fields
}

// "name": {{STRING}},
func (p *Parser) parseStructField() *ast.StructField {

	structField := &ast.StructField{
		Tok:  p.nextToken,
		Name: "",
	}

	// 字段名称
	p.NextTokenMust(token.STR)
	name := getStringVal(p.curToken.Lit)
	p.isIdentifier(name)
	structField.Name = name

	// :
	p.NextTokenMust(token.COLON)
	// 类型
	structField.Type = p.parseFieldType()

	// 检查是否有逗号, 如果有就跳过
	if p.nextToken.Tok == token.COMMA {
		p.NextTokenMust(token.COMMA)
	}
	return structField
}

// 解析字符表示的类型，例如："{{STRING}}"
func (p *Parser) parseFieldTypeFromStr() ast.Type {
	p.NextTokenMust(token.STR)
	// 从字符串中解析类型
	src := []byte(getStringVal(p.curToken.Lit))
	pos := scanner.Pos{
		File: p.curToken.Pos.GetFile(),
		Row:  p.curToken.Pos.GetRow(),
		Col:  p.curToken.Pos.GetCol() - len(p.nextToken.Lit),
	}
	e := p.parseInsertBytes(src, pos, p.parseFieldTypeWrapper)
	if e != nil {
		panic(e)
	}
	return p.typ
}

func (p *Parser) parseFieldTypeWrapper() error {
	// 存储在全局临时空间中，上层使用
	p.typ = p.parseFieldType()
	return nil
}

// 解析含有{{}}的type
func (p *Parser) parseFieldTypeHasBrace() ast.Type {
	p.NextTokenMust(token.LBRACE)
	p.NextTokenMust(token.LBRACE)
	t := p.parseFieldType()
	p.NextTokenMust(token.RBRACE)
	p.NextTokenMust(token.RBRACE)
	return t
}

func (p *Parser) parseFieldType() ast.Type {
	switch p.nextToken.Tok {
	case token.INT, token.INT32, token.INT64,
		token.UINT32, token.UINT64, token.DOUBLE, token.FLOAT, token.BOOL, token.BOOLEAN, token.STRING:
		p.Scan()
		return ast.NewBasicType(p.curToken)
	case token.MAP:
		return p.parseMap()
	case token.ENUM:
		return p.parseEnumType()
	case token.LBRACK:
		return p.parseArray()
	case token.LBRACE:
		return p.parseFieldTypeHasBrace()
	case token.STR:
		return p.parseFieldTypeFromStr()
	case token.IDENT:
		p.NextTokenMust(token.IDENT)


		return ast.NewStructType(p.curToken)
	default:
		p.errExpect("field type")
	}
	return nil
}

func (p *Parser) parseArray() *ast.Array {
	array := ast.NewArrayType(p.nextToken)
	p.NextTokenMust(token.LBRACK)
	array.ItemType = p.parseFieldType()
	p.NextTokenMust(token.RBRACK)
	return array
}

func (p *Parser) parseMap() *ast.Map {
	m := ast.NewMapType(p.nextToken)
	p.NextTokenMust(token.MAP)
	p.NextTokenMust(token.LSS)
	m.Key = p.parseMapKeyType()
	p.NextTokenMust(token.COMMA)
	m.Val = p.parseMapValType()
	p.NextTokenMust(token.GTR)
	return m
}

func (p *Parser) parseMapKeyType() ast.Type {
	return p.parseFieldType()
}

func (p *Parser) parseMapValType() ast.Type {
	return p.parseFieldType()
}

func (p *Parser) parseEnumType() ast.Type {
	p.NextTokenMust(token.ENUM)
	enum := ast.NewEnumType(p.curToken)

	// 兼容老版本，如果没有<>
	// 当做STRING类型处理
	if p.nextToken.Tok != token.LSS {
		p.curToken.Tok = token.STRING
		return ast.NewBasicType(p.curToken)
	}
	p.NextTokenMust(token.LSS)

	enum.ValType = p.parseFieldType()
	p.NextTokenMust(token.GTR)
	return enum
}

/*
* enum "city" : {
*	"beijin":  1",
*   "shanghai": 2,
* }
 */
func (p *Parser) parseEnum() {
	enum := &ast.Enum{
		Tok:  p.nextToken,
		Name: "",
	}

	p.NextTokenMust(token.ENUM)
	tokInfo := p.NextTokenMust(token.STR)

	enum.Name = getStringVal(tokInfo.Lit)
	p.isIdentifier(enum.Name)

	p.NextTokenMust(token.COLON)
	p.NextTokenMust(token.LBRACE)

	if p.nextToken.Tok != token.STR {
		p.errorInfo("enum value can not be empty")
	}

	for p.nextToken.Tok != token.RBRACE {
		enum.Fields = append(enum.Fields, p.parseEnumField())
	}
	p.NextTokenMust(token.RBRACE)

	p.AddEnumToPool(enum)
}

func (p *Parser) parseEnumField() *ast.EnumField {
	field := &ast.EnumField{
		Tok: p.nextToken,
	}
	tokInfo := p.NextTokenMust(token.STR)
	field.Name = getStringVal(tokInfo.Lit)
	p.isIdentifier(field.Name)
	p.NextTokenMust(token.COLON)
	tokInfo = p.NextTokenMust(token.INTEGER)
	tag, _ := strconv.Atoi(tokInfo.Lit)
	field.Tag = ast.EnumTag{
		Tag: tag,
		Tok: tokInfo,
	}
	// 检查是否有逗号, 如果有就跳过
	if p.nextToken.Tok == token.COMMA {
		p.NextTokenMust(token.COMMA)
	}
	return field
}

func (p *Parser) parseRequest() {

	req := &ast.Request{}

	p.NextTokenMustWithLit(token.STR, `"request"`)
	p.NextTokenMust(token.COLON)
	p.NextTokenMust(token.LBRACE)
	p.NextTokenMustWithLit(token.STR, `"url"`)
	p.NextTokenMust(token.COLON)

	tokenInfo := p.NextTokenMust(token.STR)
	req.Url = ast.RequestUrl{
		Tok:  tokenInfo,
		Name: getStringVal(tokenInfo.Lit),
	}

	p.AddRequestToPool(req)
	if p.nextToken.Tok == token.COMMA {
		p.NextTokenMust(token.COMMA)
	}

	// 如果下一个直接是'}',说明没有params,
	if p.nextToken.Tok == token.RBRACE {
		p.NextTokenMust(token.RBRACE)
		return
	}

	p.NextTokenMustWithLit(token.STR, `"params"`)
	paramToken := p.curToken

	p.NextTokenMust(token.COLON)

	// 如果是 { 根据规则自动生成struct
	if p.nextToken.Tok == token.LBRACE {
		fields := p.parseStructBody()
		if len(fields) == 0 {
			// 如果为空，则认为没有params
			return
		}
		paramStructName := req.GeneParamsName()
		// 添加param
		strcu := &ast.Struct{
			Tok:    paramToken,
			Name:   paramStructName,
			Fields: fields,
		}
		p.AddStructToPool(strcu)
		req.Params = &ast.RequestParams{
			Tok:        p.curToken,
			StructName: paramStructName,
		}
	} else {
		// 如果不是{，那么就应该是直接指定了struct的名称
		p.parseFieldType()
		if p.nextToken.Tok == token.COMMA {
			p.NextTokenMust(token.COMMA)
		}
		typ, ok := p.typ.(ast.StructType)
		if !ok {
			p.errorInfo("params must be a struct type")
		}
		req.Params = &ast.RequestParams{
			Tok:        paramToken,
			StructName: typ.Name,
		}
	}
	p.NextTokenMust(token.RBRACE)
}

func (p *Parser) parseResource() {

}

func (p *Parser) parsePackage() {
	p.NextTokenMust(token.PACKAGE)
	p.NextTokenMust(token.IDENT)
	p.NextTokenMust(token.SEMICOLON)
}

func (p *Parser) parseImport() {
	p.NextTokenMustWithLit(token.IDENT, "import")
	tokenInfo := p.NextTokenMust(token.STR)
	fileName := getStringVal(tokenInfo.Lit)
	p.NextTokenMust(token.SEMICOLON)

	// curFileName
	curFileName := p.scan[p.level].FileName()
	curImp := p.getImport(curFileName)

	// parse import file
	// 检查全局是否已经import过了
	hasImported := p.isImported(fileName)

	// 检查当前文件是否已经import过当前文件
	// 重复引用报错
	if hasImported {
		for _, n := range curImp.List {
			if n.Name == fileName {
				p.errorInfo(fileName + " has been imported")
				return
			}
		}
	}

	nextImp := p.getImport(fileName)
	nextImp.InCount++
	curImp.List = append(curImp.List, nextImp)

	// 检查是否有循环import
	if ok, path := p.isRecursiveImport(); ok {
		panic(err.Err("File recursively imports itself: " + path))
	}

	if hasImported {
		return
	}
	e := p.ParseFile(fileName)
	if e != nil {
		panic(e)
	}
}

func (p *Parser) Scan() *scanner.TokenInfo {
	p.curToken = p.nextToken
	p.nextToken = p.scan[p.level].Scan()

	if p.curToken != nil {
		p.appendToken(p.curToken)
	}
	for p.nextToken.Tok == token.EMPTY || p.nextToken.Tok == token.COMMENT {
		p.appendToken(p.nextToken)
		p.nextToken = p.scan[p.level].Scan()
	}
	return p.curToken
}

func (p *Parser) appendToken(tokenInfo *scanner.TokenInfo) {
	fmt.Print(tokenInfo.Lit)
	p.tokens = append(p.tokens, tokenInfo)
}

func (p *Parser) NextTokenMust(tok token.Token) *scanner.TokenInfo {
	tokenInfo := p.Scan()
	if tokenInfo.Tok != tok {
		p.errExpectBut(tok.String(), p.curToken.Tok.String())
	}
	return tokenInfo
}

func (p *Parser) NextTokenMustWithLit(tok token.Token, lit string) *scanner.TokenInfo {
	tokenInfo := p.NextTokenMust(tok)
	if tokenInfo.Lit != lit {
		p.errExpectBut(lit, tokenInfo.Lit)
	}
	return tokenInfo
}

func (p *Parser) ParseFile(fileName string) error {
	var src []byte
	var e error

	if len(p.protoPath) == 0 {
		p.protoPath = append(p.protoPath, "./")
	}
	for _, path := range p.protoPath {
		path = strings.Trim(path, " ")
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		file := path + fileName
		src, e = ioutil.ReadFile(file)
		if e != nil {
			continue
		}
		break
	}

	if e != nil {
		return e
	}
	pos := scanner.Pos{
		File: fileName,
		Row:  1,
	}
	return p.parseInsertBytes(src, pos, p.Start)
}

func (p *Parser) parseInsertBytes(src []byte, pos scanner.Pos, run func() error) error {
	var e error
	defer func() {
		if er := recover(); er != nil {
			switch er.(type) {
			case error:
				e = er.(error)
			default:
				e = fmt.Errorf("%v", er)
			}
		}
	}()
	s := &scanner.Scanner{}
	s.Init(src, pos.GetFile(), pos.GetRow(), pos.GetCol())
	p.scan = append(p.scan, s)
	p.level = len(p.scan) - 1

	// 保留现场
	curToken := p.curToken
	nextToken := p.nextToken
	p.tokenPointInit()
	e = run()
	if e != nil {
		return e
	}
	p.scan = p.scan[:p.level]
	p.level--
	// 恢复现场
	p.curToken = curToken
	p.nextToken = nextToken
	return nil
}

// 初始化token的指针
// 保持curToken = nil && nextToken != nil
func (p *Parser) tokenPointInit() {
	// 判断是否已经处于初始化状态
	if p.curToken == nil && p.nextToken != nil {
		return
	}
	// 不出与初始化状态
	p.curToken = nil
	p.nextToken = nil
	p.Scan()
}

func (p *Parser) AddStructToPool(stru *ast.Struct) {
	if p.enumPools == nil {
		p.enumPools = make(map[string]*ast.Enum)
	}
	if p.structPools == nil {
		p.structPools = make(map[string]*ast.Struct)
	}
	// check struct 是否已经存在
	if r, ok := p.structPools[stru.Name]; ok {
		pos := r.Tok.Pos.String()
		p.errorInfo(fmt.Sprintf("struct:'%s' already has existed on %s", r.Name, pos[:len(pos)-1]))
	}

	// check enum 是否已经存在
	if r, ok := p.enumPools[stru.Name]; ok {
		pos := r.Tok.Pos.String()
		p.errorInfo(fmt.Sprintf("enum:'%s' already has existed on %s", r.Name, pos[:len(pos)-1]))
	}

	// 放入池中
	p.structPools[stru.Name] = stru
}

func (p *Parser) AddRequestToPool(req *ast.Request) {
	if p.requestPool == nil {
		p.requestPool = make(map[string]*ast.Request)
	}
	// check url 是否已经存在
	if r, ok := p.requestPool[req.Url.Name]; ok {
		pos := r.Url.Tok.Pos.String()
		p.errorInfo(fmt.Sprintf("url:'%s' already has existed on %s", r.Url.Name, pos[:len(pos)-1]))
	}
	// 放入池中
	p.requestPool[req.Url.Name] = req
}

func (p *Parser) AddEnumToPool(enum *ast.Enum) {
	if p.enumPools == nil {
		p.enumPools = make(map[string]*ast.Enum)
	}
	if p.structPools == nil {
		p.structPools = make(map[string]*ast.Struct)
	}
	// check struct 是否已经存在
	if r, ok := p.structPools[enum.Name]; ok {
		pos := r.Tok.Pos.String()
		p.errorInfo(fmt.Sprintf("struct:'%s' already has existed on %s", r.Name, pos[:len(pos)-1]))
	}

	// check enum 是否已经存在
	if r, ok := p.enumPools[enum.Name]; ok {
		pos := r.Tok.Pos.String()
		p.errorInfo(fmt.Sprintf("enum:'%s' already has existed on %s", r.Name, pos[:len(pos)-1]))
	}

	// 放入池中
	p.enumPools[enum.Name] = enum
}

func (p *Parser) AddResourceToPool(resource *ast.Resource) {
	if p.resourcePool == nil {
		p.resourcePool = make(map[string]*ast.Resource)
	}
	// check url 是否已经存在
	if r, ok := p.resourcePool[resource.Key()]; ok {
		pos := r.Header.Tok.Pos.String()
		p.errorInfo(fmt.Sprintf("resource namespace:'%s', name:'%s' already has existed on %s", r.Header.Namespace, r.Header.Name, pos[:len(pos)-1]))
	}
	// 放入池中
	p.resourcePool[resource.Key()] = resource
}