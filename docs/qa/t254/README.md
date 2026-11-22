# QA t254: só o botão, uma seta, sem confirmação

Pedido do Erick: não gostou do rodapé nem das duas setas (`[ SETUP GUIDE ↗ ]↗`), e pediu para tirar o Confirm button, com Quick Reply direto.

## O que mudou

1. O rodapé “React with 1–10…” saiu do card. Reagir, responder o título e `!ig tap` continuam funcionando.
2. A ponte não coloca mais ↗ em nenhum rótulo. Fica `[ SETUP GUIDE ]`. O Beeper ainda acrescenta um único ↗ depois dos links externos.
3. O GET no link de tap dispara o quick reply ou o postback e devolve uma página mínima “Enviado ✓” (e tenta fechar a aba). Não há mais formulário de confirmação.

Idempotência: um envio bem-sucedido do mesmo token não se repete enquanto a ponte estiver no ar (memória do processo, com trava para cliques simultâneos). Um envio que falhou pode ser tentado de novo. Previews e crawlers (Slackbot, facebookexternalhit, WhatsApp, TelegramBot, Twitterbot, Discordbot, Googlebot, curl, HEAD, `Purpose: prefetch`, `Sec-Purpose`) recebem 204 e não disparam. Cada envio, duplicata e preview ignorado entra no log, sem o token assinado.

Nenhum link de tap foi clicado nesta sessão: agora o clique envia de verdade ao Instagram. O handler foi coberto com testes e um mock do client.

## Refresh-all

22/11/2026, 11:22 America/Sao_Paulo, 42 portais: **updated 81, skipped 440, failed 0**. Krish: updated 3, skipped 4, failed 0. Nenhum SendPostback.

O HTML conferido depois do refresh (Liam e Krish) é `[ SETUP GUIDE ]`, `[ GITHUB REPO ]`, `[ YES ]`, `[ NO THANKS ]`, `[ Following ]`, `[ Send me the link ]`, sem rodapé e sem ↗ da ponte.

## Prints

Os “antes” são prints ao vivo do Beeper, no layout anterior. Os “depois” de card são reconstruções fiéis desse HTML, porque o painel da conversa no Beeper desta sessão ficou em branco (o foco automático não pintou o chat). Não são prints ao vivo do depois.

- [Antes — duas setas](01-antes-duas-setas.png)
- [Antes — rodapé](02-antes-rodape.png)
- [Depois — links](03-depois-links.png)
- [Depois — botões](04-depois-botoes.png)
- [Página Enviado ✓](05-enviado.png)
- [Antes e depois — links](06-antes-depois-links.png)
- [Antes e depois — botões](07-antes-depois-botoes.png)
- [Vídeo](qa-t254.mp4)

Commits: `0a8996c`, `bd34f56`, `19d6ba0`. Binário no ar: sha256 `3d977239b9afa3a4591f356c88ba4f318635068dda2fc7799694b9f8083fbe7d`. Backup: `mautrix-instagram.03fa5df`.
