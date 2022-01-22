create table if not exists pharmacy (
  code char(5),
  name varchar(50),
  addr_line_1 varchar(50),
  addr_line_2 varchar(50),
  addr_line_3 varchar(50),
  addr_line_4 varchar(50),
  postcode varchar(10),
  lat decimal,
  lng decimal
);
