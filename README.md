# Things to fix/add
* Backend
  * update search
  * refactor adminfunc
  * protect configTx POST methods
  * add IDs for configTx POST methods (constants)
  * better error handling (not panic)
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
  * User enters type and payload of transaction
  * Data is sent to REST and returns Hash
  * User enters private key to sign Hash
  * Signature and Hash is sent to REST
* Database Component
