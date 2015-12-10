### ebnf грамматика для o.t.

*ebnf-схема содержит специфичные для `golang.org/x/exp/ebnf` символы экранирования и т.д.*

Символ `Object` является стартовым. Символ `unicode letter` не описан, по смыслу он обозначает видимые символы юникода, то что определяется как `golang/unicode.isLetter()`. Требует уточнения.

````ebnf
Object = Qid [(":" | "::") Content {WhiteSpace Content} ";"].

Qid = [Tpl "~" ] Cls ["(" Id ")"].

Content = Object | Value.

Value = Trilean | Number | Inf | String | Object | Character.

String = SubString [{":" SubString}].

SubString = (Quote {"unicode letter"} Quote) | Character.

Character = digit {hexDigit} "U".

Quote = "`" | "'" | `"`.

Inf = ["-"]"inf".

Trilean = "true" | "false" | "null".

Number = ["-"] digit {digit} (["." digit | {digit}] | {hexDigit} "H").

Tpl = ident.

Cls = ident.

Id = ident.

ident = "unicode letter" | "$" [digit | "unicode letter" | "@" | "#" | "$" | "%" | "^" | "&" | "*" | "-" | "_" | "=" | "+" | "," | "." | "?" | "!" | "/" | "|" | "\"].

digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9".

hexDigit = "A" | "B" | "C" | "D" | "E" | "F" | digit.

space = " " | "\n" | "\r" | "\t".

WhiteSpace = space {space}.

````
