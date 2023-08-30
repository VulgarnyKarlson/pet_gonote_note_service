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
