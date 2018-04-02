package money

import "fmt"
import "encoding/xml"

type Money int

func (m Money) String() string {
  f := float32(m)
  return fmt.Sprintf("%.2f", f/100)
}

func (m Money) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
  return e.EncodeElement(m.String(), start)
}
