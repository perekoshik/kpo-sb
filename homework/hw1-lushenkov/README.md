# Zoo (HW #1 — SOLID и DI)

Консольное приложение для учёта животных и вещей Московского зоопарка.  


## Быстрый старт

```bash
go run ./cmd/zoo
```

### Примеры команд
```
add thing table BossTable
add thing computer PC-42
add animal rabbit "Bunny" 3 8
add animal tiger "ShereKhan" 30
report balance
report food
report interactive
help
exit
```

---

## применение SOLID

- **S (Single Responsibility)**  
  Каждый компонент отвечает за одну задачу:  
  - `onboarding.Service` - приём животных (ветосмотр и запись).  
  - `reports.*` - генерация отчётов.  
  - `vetclinic.SimpleClinic` - определение состояния здоровья животного.

- **O (Open/Closed)**  
  Новые виды животных или клиник можно добавлять через новые типы/реализации без изменения существующего кода.

- **L (Liskov Substitution)**  
  Все животные подставимы как `AnimalEntity` или `Alive`.  
  Травоядные реализуют дополнительный интерфейс `KindnessCarrier`.

- **I (Interface Segregation)**  
  Небольшие интерфейсы:  
  `animal.Alive`, `inventory.InventoryItem`, `animal.KindnessCarrier`,  
  `ports.{AnimalRepository,ThingRepository,InventoryQuery}`, `vetclinic.Clinic`.

- **D (Dependency Inversion)**  
  Реализации подставляются DI-контейнером `dig` в `infra/di`.

---

## Свойства животных

- **Травоядные** имеют параметр `Kindness` (0..10).  
  Животное попадает в контактный зоопарк, если `kindness > 5`.
- **Хищники** не имеют доброты и не участвуют в контактном зоопарке.
- **Ветклиника** решает, принимать ли животное, исходя из здоровья

---

## Тестирование

```bash
go test ./... -cover
```
Общие тесты покрыты > 60 %
