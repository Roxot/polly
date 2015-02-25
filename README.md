###POST api.polly.com/register:###

* Request body:
    * [phone_number] a valid (dutch) phone number
* Returns 200 OK for success
    * Returns 400 BAD REQUEST when not providing a valid phone number
    * Returns 500 INTERNAL SERVER ERROR if one occurs

###POST api.polly.com/register/verify:###
* Request body:
    * [phone_number] a valid (dutch) phone number
    * [verification_token] the verification token (for always 'VERIFY')
    * [display_name] the user's display name
    * [device_type] "0" for android, "1" for ios
    * [device_guid] the device guid for push messages, ignored for now
* Returns 200 OK for success, containing the user in JSON in the response body:
```
#!json
        {
            "id" : 0,
            "phone_number" : "0612345678",
            "display_name" : "The wolf",
            "device_type" : 0,
            "token" : "9d8592a7-585c-438b-9db0-29465cd66c25"
        }
```
* Returns 400 BAD REQUEST when not providing all values or providing bad values
* Returns 500 INTERNAL SERVER ERROR if one occurs

###POST api.polly.com/poll:###
* Requires the use of BasicAuth using the phone number and token as username and password.
* Request body should contain pure JSON of the form:
```
#!json
        {
            "meta_data" : {
                "title" : "Filmpje doen"
            },

            "questions" : [
                ...,

                {
                    "id" : 2,
                    "type" : 0,                                         <-- 0 = multiple_choice, 1 = open, 2 = date
                    "title" : "Naar welke film gaan we vanavond?"
                },

                ...
            ],

            "options" : [
                {
                    "id" : 0,
                    "question_id" : 2,
                    "value" : "The imitation game"
                },

                ...
            ],

            "participants" : [                                          <-- Should not contain self
                {
                    "id" : 10298
                    "phone_number" : "0687654321"
                    "display_name" : "Friend of Polly Client"
                },

                ...
            ]

        }
```
* Returns 200 OK for success, containing the poll in JSON in the response body:
```
#!json        
        {
            "meta_data" : {
                "poll_id" : 283                                              <-- Server-side id
                "creation_date" : "1073029382"                          <-- Unix time
                "title" : "Filmpje doen"
            },

            "questions" : [
                ...,

                {
                    "id" : 1231,                                        <-- Server-side id
                    "type" : 0,
                    "title" : "Naar welke film gaan we vanavond?"
                },

                ...
            ],

            "options" : [
                {
                    "id" : 1923,                                        <-- Server-side id
                    "question_id" : 2,
                    "value" : "The imitation game"
                },

                ...
            ],

            "creator" : {
                "id" : 1073,
                "phone_number" : "0612345678",
                "display_name" : "Polly Client"
            },

            "votes" : [],                                               <-- Always empty, can be ignored

            "participants" : [
                {
                    "id" : 10298
                    "phone_number" : "0687654321"
                    "display_name" : "Friend of Polly Client"
                },

                ...
            ]

        }
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.
* Returns 500 INTERNAL SERVER ERROR if one occurs

###GET api.polly.com/poll/xx:###
* Requires the use of BasicAuth using the phone number and token as username and password.
* Replace xx with the server-side identifier of the poll you're requesting.
* Returns 200 OK for success, containing the poll in JSON in the response body:
```
#!json
        {
            "meta_data" : {
                "id" : 283
                "creation_date" : "1073029382"
                "title" : "Filmpje doen"
            },

            "questions" : [
                ...,

                {
                    "id" : 1231,
                    "type" : 0,
                    "title" : "Naar welke film gaan we vanavond?"
                },

                ...
            ],

            "options" : [
                {
                    "id" : 1923,
                    "question_id" : 2,
                    "value" : "The imitation game"
                },

                ...
            ],

            "creator" : {
                "id" : 1073,                                                            <-- user id
                "phone_number" : "0612345678",
                "display_name" : "Polly Client"
            },

            "votes" : [
                {
                    "id" : 102309
                    "option_id" : 1923
                    "user_id" : 1073
                },

                ...
            ],

            "participants" : [
                {
                    "id" : 10298                                                        <-- user id
                    "phone_number" : "0687654321"
                    "display_name" : "Friend of Polly Client"
                },

                ...
            ]

        }
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.