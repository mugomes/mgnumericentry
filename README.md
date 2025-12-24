# MGNumericEntry

MGNumericEntry é um componente numérico para Fyne.

## Recursos

- Limite o valor numérico (mínimo e máximo)
- Adicione um valor inicial
- Botão para incrementar e decrementar

## Instalação

`go get github.com/mugomes/mgnumericentry`

## Exemplo

```
import "github.com/mugomes/mgnumericentry"

campoNumerico, retorno := mgnumericentry.NewMGNumericEntryWithButtons(0, 100, 20)
print(retorno.GetValue())
```

## Information

 - [Page MGNumericEntry](https://github.com/mugomes/mgnumericentry)

## Requirement

 - Go 1.24.6
 - Fyne 2.7.0

## Support

- GitHub: https://github.com/sponsors/mugomes
- More: https://mugomes.github.io/apoie.html

## License

Copyright (c) 2025 Murilo Gomes Julio

Licensed under the [MIT](https://github.com/mugomes/mgnumericentry/blob/main/LICENSE) license.

All contributions to the MGNumericEntry are subject to this license.
