Atualização t254 (só o botão, sem rodapé, sem ↗ da ponte, sem confirmação): ver [t254/README.md](t254/README.md) e o [vídeo novo](t254/qa-t254.mp4).

# QA visual: pseudo-botões `[ Título ]` no Beeper

O Beeper não tem UI de botão. A ponte passa a renderizar cada CTA (quick reply, postback e link) como um pseudo-botão, um por linha:

```
[ BOTÃO 1 ]
[ BOTÃO 2 ]
```

O texto inteiro `[ Título ]` é o link clicável. Há um espaço inseparável (NBSP) entre o colchete e o título, para o controle não quebrar. As linhas são unidas por `<br>`, sem `<ul>`/`<li>`. O ↗ de link externo fica colado dentro dos colchetes (`[ Título ↗ ]`). Não há instrução longa em cada botão. Só um rodapé discreto com os jeitos alternativos (reação, título ou `!ig tap N`), e apenas quando o card tem quick reply ou postback.

Commit no branch `ig-xma-postback-buttons`: `03fa5df` (`igconv: render tappable CTAs as bracketed pseudo-buttons`).

## Prints

1. [Antes — o ↗ cai sozinho](01-antes-quebra.png)
2. [Antes — texto longo em cada botão](02-antes-lista.png)
3. [Antes — o mesmo problema no Krish](03-antes-krish-links.png)
4. [Depois — links em pseudo-botão](04-depois-links.png)
5. [Depois — [ YES ] e [ NO THANKS ]](05-depois-botoes.png)
6. [Depois — um único botão tocável](06-depois-send.png)
7. [Antes e depois — links numa janela estreita](07-antes-depois-links.png)
8. [Antes e depois — botões tocáveis](08-antes-depois-botoes.png)
9. [Durante — página de confirmação, token omitido](09-durante-confirmacao.png)

Há também um [vídeo curto](qa-pseudo-botoes.mp4) com esses frames.

Os prints de depois foram tirados no Beeper depois do `!ig refresh-all`. Recortes da conversa, sem a lista da inbox. Nomes de terceiros não aparecem nesses recortes. O token assinado da confirmação foi omitido.

## Atualização das mensagens antigas

`!ig refresh-all` na sala de gestão, em 22/11/2026, 10:49 America/Sao_Paulo:

- Totais em 42 portais: **updated 81, skipped 437, failed 0**
- Krish (`18093156415972515`): updated 3, skipped 4, failed 0
- Nada foi enviado ao Instagram (sem Send, sem POST, sem SendPostback)

As mensagens novas já nascem nesse formato. As antigas foram refeitas a partir dos dados que a ponte já guarda.

## O que ainda é limitação do Beeper

O Beeper acrescenta um ↗ depois de qualquer `<a>`, inclusive depois do `]` dos pseudo-botões, e também nos botões tocáveis (a ponte não coloca ↗ nesses). O ↗ de dentro do colchete, só em links externos, permanece colado no título. Neste card estreito, o ↗ extra do Beeper não caiu sozinho na linha de baixo.

O clique hoje aponta para `http://127.0.0.1:29330/_ig/tap/...`. No computador da ponte, o GET só abre a confirmação; o POST (Send) é que dispara o toque. No Mac e no celular o clique ainda não chega na ponte, até existir um endereço público, com autorização explícita.

Nada neste material deve ir para o PR oficial [mautrix/meta#639](https://github.com/mautrix/meta/pull/639).
