package parser



type Import struct {
	Name string
	// 入度
	InCount int

	// next list
	List []*Import
}


// 检查是否有循环import
func (p *Parser) isRecursiveImport() (bool, string) {

	if len(p.imports) <= 1 {
		return false, ""
	}

	m := make(map[string]*Import, len(p.imports))
	for k, v := range p.imports {
		vv := *v
		m[k] = &vv
	}
	// 深度copy Import.List
	for k, v := range p.imports {
		m[k].List = make([]*Import, len(v.List))
		for i, n := range  v.List {
			m[k].List[i] = m[n.Name]
		}
	}

	// 判断是否有环
	flag := true
	for flag {
		flag = false
		for k, node := range m {
			if node.InCount == 0 {
				for _, next := range node.List {
					next.InCount--
				}
				delete(m, k)
				flag = true
				break
			}
		}
	}

	if len(m) > 0 {
		path := ""
		var node *Import
		for _, node = range m {
			break
		}
		if node == nil {
			return false, ""
		}
		startName := node.Name
		path = node.Name
		for {
			node = node.List[0]
			path += "--> " + node.Name
			if node.Name == startName {
				break
			}
		}
		return true, path
	}
	return false, ""
}


func (p *Parser) isImported(file string) bool{
	_, ok := p.imports[file]
	return ok
}
func (p *Parser) getImport(file string)*Import {
	if p.imports == nil {
		p.imports = make(map[string]*Import)
	}
	imp := p.imports[file]
	if imp == nil {
		imp = &Import{
			Name:    file,
		}
		p.imports[file] = imp
	}
	return imp
}
