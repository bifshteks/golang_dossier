# Go
Чтобы начать нормально работать с новым проектом
```bash
go mod init
```

что создаст в текущей папке go.mod, типо requirements.txt, с названием модуля.
После того, как в этой папке появился этот go.mod, можно устанавливать зависимости через
```bash
go get <packageName>
```
которые запишутся в зависимости этого пакета