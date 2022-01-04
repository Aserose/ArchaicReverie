# ArchaicReverie

На данный момент **«ArchaicReverie»** - это пока только гейм-концепт (набросок, задача которого показать суть игровой механики). Главной задумкой является реализация системы, где у каждого пользователся есть возможность создавать своих персонажей и затем выбирать одного из них для решения задач, сгенерированных условиями игрового события (куда входят характеристики местности, времени суток, погоды). Например, прыжок: игрок задаёт положение самого прыжка (амплитуда движения рук, разбег и т.д.), потом код рассчитывает по специальной формуле все параметры с последующей проверкой на успешность действия (которая, в свою очередь, является результатом вычислений по собственной формуле упомянутых характеристик игрового события) и отображает итог. 
___
### Конструкция

API разработана с применением фреймворка Gin.

Авторизация происходит через токен и куки. 

Репозиторий: SQL СУБД PostgreSQL

Среди задействованных библиотек имеется: 
* sqlx (для работы с SQL СУБД PostgreSQL);
* gin (API);
* jwt (токен)

и так далее. 

**Сделано**:

#04.01.2022
- генерация игровых условий;
- логика прыжка;
- тестирование имеющегося;
- мелкие доработки;

#...
- система авторизации;
- система управления профилями персонажей через аккаунт пользователя (создание, выбор, редактирование, удаление);

**TODO**:
- [ ] логика:
	- [ ] игровых условий
	- [ ] действия персонажа;
- [ ] взаимодействие с другими персонажами;
- [ ] механика генерации случайных задач на совместную работу;
- [ ] фронтенд;
