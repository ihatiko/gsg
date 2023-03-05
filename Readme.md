Что хочу видеть ?

1) Добавить опциональную очистку базы
2) Язык правил
3) Intellisence для этих правил (vs code)
4) Аннотации для связи между базами
5) Контексты для связи между базами
6) Управление кол-во записей через конфиг
7) Работает как golang приложение через go install
8) добавить смещение serial id (nextval)
9) Прописать все constrains
10) Работа со справочниками
11) Прекомпиляция и валидация правил на справочники
12) Исправить конфликты при генерации уникальных значений
13) Преобразователь типов между базами
14) Добавить правило срабатывание null в проценте
15) Добавить срабатывание валидации на кол-ве сетов и возможных unique значений


create table test (
id serial
);

1. alter table test add column case0 **serial**;
2. alter table test add column case1 **uuid**;
3. alter table test add column case2 **bit**;
4. alter table test add column case3 **bool**;
5. alter table test add column case4 **date**;
6. alter table test add column case5 **timestamp**;
7. alter table test add column case6 **numeric**;
8. alter table test add column case7 **varchar(256)**;
9. alter table test add column case7 **jsonb**;
10. alter table test add column case7 **text**;

Карта поддерживаемых типов
Name	Aliases	Description Supported
1. [ ] bigint	int8	signed eight-byte integer
2. [ ] bigserial	serial8	autoincrementing eight-byte integer
3. [ ] bit [ (n) ]	 	fixed-length bit string
4. [X] bit varying [ (n) ]	varbit [ (n) ]	variable-length bit string
5. [X] boolean	bool	logical Boolean (true/false)
6. [ ] box	 	rectangular box on a plane
7. [ ] bytea	 	binary data (“byte array”)
8. [ ] character [ (n) ]	char [ (n) ]	fixed-length character string
9. [ ] character varying [ (n) ]	varchar [ (n) ]	variable-length character string
10. [ ] cidr	 	IPv4 or IPv6 network address
11. [ ] circle	 	circle on a plane
12. [X] date	 	calendar date (year, month, day)
13. [ ] double precision	float8	double precision floating-point number (8 bytes)
14. [ ] inet	 	IPv4 or IPv6 host address
15. [ ] integer	int, int4	signed four-byte integer
16. [ ] interval [ fields ] [ (p) ]	 	time span
17. [ ] json	 	textual JSON data
18. [ ] jsonb	 	binary JSON data, decomposed
19. [ ] line	 	infinite line on a plane
20. [ ] lseg	 	line segment on a plane
21. [ ] macaddr	 	MAC (Media Access Control) address
22. [ ] macaddr8	 	MAC (Media Access Control) address (EUI-64 format)
23. [ ] money	 	currency amount
24. [X] numeric [ (p, s) ]	decimal [ (p, s) ]	exact numeric of selectable precision
25. [ ] path	 	geometric path on a plane
26. [ ] pg_lsn	 	PostgreSQL Log Serial Number
27. [ ] pg_snapshot	 	user-level transaction ID snapshot
28. [ ] point	 	geometric point on a plane
29. [ ] polygon	 	closed geometric path on a plane
30. [ ] real	float4	single precision floating-point number (4 bytes)
31. [ ] smallint	int2	signed two-byte integer
32. [ ] smallserial	serial2	autoincrementing two-byte integer
33. [X] serial	serial4	autoincrementing four-byte integer
34. [ ] text	 	variable-length character string
35. [ ] time [ (p) ] [ without time zone ]	 	time of day (no time zone)
36. [ ] time [ (p) ] with time zone	timetz	time of day, including time zone
37. [X] timestamp [ (p) ] [ without time zone ]	 	date and time (no time zone)
38. [ ] timestamp [ (p) ] with time zone	timestamptz	date and time, including time zone
39. [ ] tsquery	 	text search query
40. [ ] tsvector	 	text search document
41. [ ] txid_snapshot	 	user-level transaction ID snapshot (deprecated; see pg_snapshot)
42. [X] uuid	 	universally unique identifier
43. [ ] xml	 	XML data
