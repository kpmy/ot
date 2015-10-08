# ot
[![Build Status](https://drone.io/github.com/kpmy/ot/status.png)](https://drone.io/github.com/kpmy/ot/latest)

Object template language

o.t. is when this
````xml
<note>
  <to>Tove</to>
  <from>Jani</from>
  <heading>Reminder</heading>
  <body>Don't forget me this weekend!</body>
</note>
````
becomes this
````
note:
  to: "Tove";
  from: "Jani";
  heading: Reminder;
  body: `Don't forget me this weekend!`;
;
````
