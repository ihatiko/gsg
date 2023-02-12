Что хочу видеть ?

1) Язык правил
2) Intellisence для этих правил (vs code)
3) Аннотации для связи между базами
4) Контексты для связи между базами
5) Управление кол-во записей через конфиг
6) Работает как golang приложение через go install
7) добавить смещение serial id (nextval)
8) Прописать все constrains
9) Работа со справочниками
10) Прекомпиляция и валидация правил на справочники


Карта поддерживаемых типов
Name	Aliases	Description Supported
1. [x] bigint	int8	signed eight-byte integer
2. [x] bigserial	serial8	autoincrementing eight-byte integer
3. [x] bit [ (n) ]	 	fixed-length bit string
4. [x] bit varying [ (n) ]	varbit [ (n) ]	variable-length bit string
5. [x] boolean	bool	logical Boolean (true/false)
6. [x] box	 	rectangular box on a plane
7. [x] bytea	 	binary data (“byte array”)
8. [x] character [ (n) ]	char [ (n) ]	fixed-length character string
9. [x] character varying [ (n) ]	varchar [ (n) ]	variable-length character string
10. [x] cidr	 	IPv4 or IPv6 network address
11. [x] circle	 	circle on a plane
12. [x] date	 	calendar date (year, month, day)
13. [x] double precision	float8	double precision floating-point number (8 bytes)
14. [x] inet	 	IPv4 or IPv6 host address
15. [x] integer	int, int4	signed four-byte integer
16. [x] interval [ fields ] [ (p) ]	 	time span
17. [x] json	 	textual JSON data
18. [x] jsonb	 	binary JSON data, decomposed
19. [x] line	 	infinite line on a plane
20. [x] lseg	 	line segment on a plane
21. [x] macaddr	 	MAC (Media Access Control) address
22. [x] macaddr8	 	MAC (Media Access Control) address (EUI-64 format)
23. [x] money	 	currency amount
24. [x] numeric [ (p, s) ]	decimal [ (p, s) ]	exact numeric of selectable precision
25. [x] path	 	geometric path on a plane
26. [x] pg_lsn	 	PostgreSQL Log Sequence Number
27. [x] pg_snapshot	 	user-level transaction ID snapshot
28. [x] point	 	geometric point on a plane
29. [x] polygon	 	closed geometric path on a plane
30. [x] real	float4	single precision floating-point number (4 bytes)
31. [x] smallint	int2	signed two-byte integer
32. [x] smallserial	serial2	autoincrementing two-byte integer
33. [x] serial	serial4	autoincrementing four-byte integer
34. [x] text	 	variable-length character string
35. [x] time [ (p) ] [ without time zone ]	 	time of day (no time zone)
36. [x] time [ (p) ] with time zone	timetz	time of day, including time zone
37. [x] timestamp [ (p) ] [ without time zone ]	 	date and time (no time zone)
38. [x] timestamp [ (p) ] with time zone	timestamptz	date and time, including time zone
39. [x] tsquery	 	text search query
40. [x] tsvector	 	text search document
41. [x] txid_snapshot	 	user-level transaction ID snapshot (deprecated; see pg_snapshot)
42. [x] uuid	 	universally unique identifier
43. [x] xml	 	XML data