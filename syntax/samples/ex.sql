/*
Example
*/

-- service providers table
CREATE TABLE dbo.service_shops
(
  [id] [int] Identity(1,1) NOT NULL,
  [name] [varchar] (50) NOT NULL,
  [onduty_code] [tinyint] NOT NULL,
  [phone] [varchar] (50)
) ON [PRIMARY];

-- this select returns 0 - if nothing found, 1 - if condition is met
select count(*)
from (select DATEPART(hour, GETDATE()) as h) as a
where a.h BETWEEN 8 AND 20;-- time between 8:00 and 20:00

-- select providers that works at this time + those that always available

select t1.name, t1.phone
from dbo.service_shops AS t1 inner join
(-- day time
  select count(*) * 1 as code
  from (select DATEPART(hour,getdate()) as h) as a
  where a.h BETWEEN 8 AND 20
) as t2 on t1.onduty_code = t2.code
    UNION ALL
select t1.name, t1.phone
from dbo.service_shops AS t1 inner join
(-- night hours
  select count(*) * 2 as code
  from (select DATEPART(hour,getdate()) as h) as a
  where (a.h BETWEEN 1 AND 7) OR (a.h BETWEEN 21 AND 24)
) as t2 on t1.onduty_code = t2.code
    UNION ALL
-- 24 hours 7 days a week
select name, phone from dbo.service_shops where onduty_code = 3; 