# object template (ot)
[![Build Status](https://drone.io/github.com/kpmy/ot/status.png)](https://drone.io/github.com/kpmy/ot/latest)

Язык описания и шаблонизации объектов ([об идее](http://b.ocsf.in/2015/10/04/%D0%A8%D0%B0%D0%B1%D0%BB%D0%BE%D0%BD%D0%B8%D0%B7%D0%B0%D1%82%D0%BE%D1%80/))

[Wiki](https://github.com/kpmy/ot/wiki)

Object description/template language

o.t. is when this
````xml
<note status="saved">
  <to>Tove</to>
  <from>Jani</from>
  <heading>Reminder</heading>
  <body>Don't forget me this weekend!</body>
</note>
````
becomes this
````
note:
  status :: saved;
  to: "Tove";
  from: "Jani";
  heading: 'Reminder';
  body: `Don't forget me this weekend!`;
;
````

## TODO

- [x] служебный блок `context` с данными пользователя
- [x] модульность шаблонов, include-фичу
- [ ] ~~схема/валидация~~ (relaxng прикрутить можно, но лень)
- [x] пути (примитивные)
- [ ] сериализация/десериализация
- [ ] ~~встроить leaf/lomo~~
- [ ] генератор текстового представления o.t.
- [ ] рбнф языка

## Итоги

Всё хорошо. 
