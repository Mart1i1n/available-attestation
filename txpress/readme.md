### build binary
```
make
```

### generate test account
create 10000 test account with eth coin.
```
./txpress account --count 1000 --eth 10000
```

then generate accounts and private key in `accounts.json`, and genesis inifo in `balance.json`.


### do test
modify app.json with your case.

- **tx_count：** total transaction count to generate and send.
- **send_routine_count：** routine count to concurrent send tx.
- **speed:** send tx count per second.
- **erc20_contract:** default used erc20 token contract.
- **rpc_node:** rpc url to the chain.
- **receive_addr：** receive address for all transaction, if empty, will create a new account to receive.
- **amount：** the amount of transfer.
- **chain_id：** chain id. 
- **batch_transfer_contract：** contract used to batch send eth/token when init test accounts.

test transfer eth coin.
```
./txpress --start
```
and it will give a report after all transaction done.
```text
INFO[2024-05-23 17:40:03.624] collect blocks txfind 100, start block is 306, endblock is 309 
INFO[2024-05-23 17:40:04.438] wait next block 310 generate                 
INFO[2024-05-23 17:40:06.194] wait next block 310 generate                 
INFO[2024-05-23 17:40:08.023] wait next block 310 generate                 
INFO[2024-05-23 17:40:09.783] wait next block 310 generate                 
INFO[2024-05-23 17:40:11.610] wait next block 310 generate                 
INFO[2024-05-23 17:40:13.451] wait next block 310 generate                 
INFO[2024-05-23 17:40:15.271] wait next block 310 generate                 
INFO[2024-05-23 17:40:17.136] block 306 have tx 28, blocktime 1716457167, cost time 12s 
INFO[2024-05-23 17:40:17.136] block 307 have tx 41, blocktime 1716457179, cost time 12s 
INFO[2024-05-23 17:40:17.136] block 308 have tx 30, blocktime 1716457191, cost time 12s 
INFO[2024-05-23 17:40:17.955] block 309 have tx 1, blocktime 1716457203, cost time 12s 
INFO[2024-05-23 17:40:18.716] total tx 100 and cost 48, tps is 2           
INFO[2024-05-23 17:40:18.716] test finished      
```