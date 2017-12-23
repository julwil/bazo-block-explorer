# Things to fix/add
* Backend
  * refactor adminfunc
  * better error handling (not panic)
  * open transactions
* JS login (client side)
  * User enters private key
  * encode JWT using private key
  * store token in Cookie
  * Save cookie
* Login (server side)
  * User sends request with Cookie
  * decode JWT using public key (from where public keys?)
  * Grant access to route
* Send TX via JS (client side)
  * make all transactions
* Database Component
