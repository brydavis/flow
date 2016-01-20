select distinct Id, FirstName into #Data
from Persons
where FirstName like 'B%'

;

select FirstName, count(Id) as N
from #Data
group by FirstName
order by count(Id) desc

;


drop table #Data