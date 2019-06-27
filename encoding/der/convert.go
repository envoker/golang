package der

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func ConvertToString(n *Node) (string, error) {
	var buf bytes.Buffer
	err := nodeToString(n, &buf, 0)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func nodeToString(n *Node, buf *bytes.Buffer, indent int) (err error) {

	indentBuff := make([]byte, indent)
	for i := 0; i < indent; i++ {
		indentBuff[i] = '\t'
	}

	if _, err = buf.Write(indentBuff); err != nil {
		return
	}

	var className string

	switch n.t.class {
	case CLASS_UNIVERSAL:
		className = "UNIVERSAL"
	case CLASS_APPLICATION:
		className = "APPLICATION"
	case CLASS_CONTEXT_SPECIFIC:
		className = "CS"
	case CLASS_PRIVATE:
		className = "PRIVATE"
	}

	s := fmt.Sprintf("%s(%d):", className, int(n.t.tagNumber))
	if _, err = buf.WriteString(s); err != nil {
		return
	}

	if n.t.valueType == VT_PRIMITIVE {

		var pPrimitive *Primitive = n.v.(*Primitive)

		buf.WriteByte('\t')

		s = hex.EncodeToString(pPrimitive.Bytes())
		if _, err = buf.WriteString(s); err != nil {
			return
		}

		buf.WriteByte('\n')

	} else if n.t.valueType == VT_CONSTRUCTED {

		buf.WriteString("\t{\n")

		var pConstructed *Constructed = n.v.(*Constructed)

		for _, child := range pConstructed.nodes {
			if err = nodeToString(child, buf, indent+1); err != nil {
				return
			}
		}

		buf.Write(indentBuff)
		buf.WriteString("}\n")
	}

	return
}
