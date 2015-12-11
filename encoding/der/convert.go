package der

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func ConvertToString(p *Node) (s string, err error) {

	var buffer bytes.Buffer

	if err = nodeToString(p, &buffer, 0); err != nil {
		return
	}

	s = string(buffer.Bytes())

	return
}

func nodeToString(p *Node, buffer *bytes.Buffer, indent int) (err error) {

	indentBuff := make([]byte, indent)
	for i := 0; i < indent; i++ {
		indentBuff[i] = '\t'
	}

	if _, err = buffer.Write(indentBuff); err != nil {
		return
	}

	var className string

	switch p.t.class {
	case CLASS_UNIVERSAL:
		className = "UNIVERSAL"
	case CLASS_APPLICATION:
		className = "APPLICATION"
	case CLASS_CONTEXT_SPECIFIC:
		className = "CS"
	case CLASS_PRIVATE:
		className = "PRIVATE"
	}

	s := fmt.Sprintf("%s(%d):", className, int(p.t.tagNumber))
	if _, err = buffer.WriteString(s); err != nil {
		return
	}

	if p.t.valueType == VT_PRIMITIVE {

		var pPrimitive *Primitive = p.v.(*Primitive)

		buffer.WriteByte('\t')

		s = hex.EncodeToString(pPrimitive.GetBytes())
		if _, err = buffer.WriteString(s); err != nil {
			return
		}

		buffer.WriteByte('\n')

	} else if p.t.valueType == VT_CONSTRUCTED {

		buffer.WriteString("\t{\n")

		var pConstructed *Constructed = p.v.(*Constructed)

		for _, node := range pConstructed.nodes {
			if err = nodeToString(node, buffer, indent+1); err != nil {
				return
			}
		}

		buffer.Write(indentBuff)
		buffer.WriteString("}\n")
	}

	return
}
