

with t as (
select urls.id,t2.device_type,t1.browser,t0.os,t4.total_clicks,urls.username,sum(case when t4.total_clicks is not null then t4.total_clicks else 0 end) from urls
left join (select id,count(*) as total_clicks from logs where logs.username='goku'  group by id order by total_clicks desc)t4 on t4.id = urls.id
left join (select t.* from (select logs.id,logs.device_type, rank() over(partition by id order by count(device_type) desc) from logs where logs.username='goku'  group by logs.id,logs.device_type) t where rank=1)t2 on urls.id = t2.id
left join (select t.* from (select logs.id,logs.browser, rank() over(partition by id order by count(browser) desc) from logs where logs.username='goku'  group by logs.id,logs.browser) t where rank=1) t1 on t2.id = t1.id
left join  (select t.* from (select logs.id,logs.os, rank() over(partition by id order by count(os) desc) from logs where logs.username='goku'  group by logs.id,logs.os) t where rank=1)t0 on t1.id = t0.id where urls.username='goku' group by urls.id,t2.device_type,t1.browser,t0.os,t4.total_clicks
order by sum desc)
select * from t 


/*
insert into logs values('http://youtube.com','goku','windows','firefox','desktop',current_timestamp,'0.0.0.0','','summer_sale','facebook','digital','qwer'),('http://youtube.com','goku','windows','firefox','desktop',current_timestamp,'0.0.0.0','','summer_sale','facebook','digital','qwer'),('http://youtube.com','goku','windows','firefox','desktop',current_timestamp,'0.0.0.0','','summer_sale','facebook','digital','qwer'),('http://youtube.com','goku','windows','firefox','desktop',current_timestamp,'0.0.0.0','','summer_sale','facebook','digital','qwer'),('http://youtube.com','goku','windows','firefox','desktop',current_timestamp,'0.0.0.0','','summer_sale','facebook','digital','qwer'),('http://youtube.com','goku','windows','firefox','desktop',current_timestamp,'0.0.0.0','','summer_sale','facebook','digital','qwer'),('http://youtube.com','goku','windows','firefox','desktop',current_timestamp,'0.0.0.0','','summer_sale','facebook','digital','qwer') ; */ 

/*
with t as (
select t4.id,t2.device_type,t1.browser,t0.os,t4.total_clicks from (select id,count(*) as total_clicks from logs  where campaign='summer_sale'  group by id order by total_clicks desc)t4
left join (select t.* from (select logs.id,logs.device_type, rank() over(partition by id order by count(device_type) desc) from logs  where campaign='summer_sale' group by logs.id,logs.device_type) t where rank=1)t2 on t4.id = t2.id
left join (select t.* from (select logs.id,logs.browser, rank() over(partition by id order by count(browser) desc) from logs  where campaign='summer_sale'  group by logs.id,logs.browser) t where rank=1) t1 on t2.id = t1.id
left join  (select t.* from (select logs.id,logs.os, rank() over(partition by id order by count(os) desc) from logs where campaign='summer_sale'  group by logs.id,logs.os) t where rank=1)t0 on t1.id = t0.id group by t4.id,t2.device_type,t1.browser,t0.os,t4.total_clicks)
select * from t;
*/

/*
select browser,browser_count,os,os_count,device_type,device_count,t0.total_clicks from (select id,count(*) as total_clicks from logs where id='rscctN' group by id)t0 full outer join  (select browser,count(*) as browser_count,row_number () over() as num from logs where id='rscctN' group by browser order by browser_count desc) t1 on 1=1
left join                                                                                                              
(select os,count(*) as os_count, row_number() over() as num from logs where id='rscctN'  group by os order by os_count desc) t2 on t1.num = t2.num                                                                                                   
left join 
(select device_type,count(*) as device_count , row_number() over() as num from logs where id='rscctN' group by device_type order by device_count desc)t3 on t2.num = t3.num
*/

