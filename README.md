# Note Service (Golang)
## Project "GoNote" - Note management system.
### Description:

### Project Workflow:
- CRUD operations for notes:
    - Creating notes.
    - Viewing notes.
    - Searching notes based on various criteria, using advanced indexes for fast data retrieval.
    - Modifying and editing notes.
    - Deleting notes.
- Every user action with notes (creation, viewing, editing, deleting)
  sends a message to RabbitMQ for further statistics aggregation.

### Code style & hardskill techniques:
- [x] Hexagonal architecture
- [x] Dependency Injection ( uber/fx )
- [x] Transactional outbox (Stats Sender is releasing stats messages from outbox table)
- [x] Stream processing ( create_note stream read, check and write to db by batch insert )
- [x] Unit tests ( testify/assert/mockgen ) ( in progress integration and e2e tests )
- [x] DDD
- [x] CQRS
- [x] Pre-commit hooks + golangci-lint


### Installation
```bash
make bin-deps
make up-build
make up-services
make migrate-up
make up
```

## License
This project is licensed under the terms of the [CC BY-NC-SA](https://creativecommons.org/licenses/by-nc-sa/4.0/legalcode) license.
