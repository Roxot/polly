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
    * [device_type] 0 for android, 1 for iOS
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
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Request body should contain JSON of the form:
```
#!json
        {
            "meta_data" : {
                "title" : "Filmpje doen"                                <-- May not be empty
            },

            "question" : {
                    "type" : 0,                                         <-- 0 = multiple_choice, 1 = open, 2 = date
                    "title" : "Naar welke film gaan we vanavond?"       <-- May not be empty
            },

            "options" : [
                {
                    "value" : "The imitation game"                      <-- May not be empty
                },

                ...
            ],

            "participants" : [                                          <-- Should contain the creator
                {
                    "id" : 10298                                        <-- The server-side user ID of the participant
                },

                ...
            ]

        }
```
* Returns 200 OK for success, containing the poll in JSON in the response body:
```
#!json        
        {
            "meta_data": {
                "poll_id": 5,                                   <-- Server-side poll ID
                "creator_id": 1,                                <-- Creator Id, thus the sender's user ID
                "creation_date": 1428946988,                    <-- Unix time
                "last_updated": 1428946988,                     <-- Unix time
                "title": "Filmpje doen"
            },
            "question": {
                "id": 5,                                        <-- Server-side question ID
                "type": 0,
                "title": "Naar welke film gaan we vanavond?"
            },
            "options": [
                {
                    "id": 29,                                    <-- Server-side option ID
                    "question_id": 5,
                    "value": "The imitation game"
                },
                {
                    "id": 30,                                    <-- Server-side option ID
                    "question_id": 5,
                    "value": "American Sniper"
                },
                {
                    "id": 31,                                    <-- Server-side option ID
                    "question_id": 5,
                    "value": "Jupiter Ascend"
                }
            ],
            "votes": [],                                         <-- Always empty
            "participants": [
                {
                    "id": 1,                                 
                    "phone_number": "0622197479",                <-- User's phone number filled in
                    "display_name": "Bryan Eikema"               <-- User's display name filled in
                }
            ]
}
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.
* Returns 401 UNAUTHORIZED if no authentication is provided.
* Returns 500 INTERNAL SERVER ERROR if one occurs

###GET api.polly.com/poll/xx:###
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Replace xx with the server-side identifier of the poll you're requesting.
* Returns 200 OK for success, containing the poll in JSON in the response body:
```
#!json
        {
            "meta_data" : {
                "id" : 283
                "creation_date" : 1073029382
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
* Returns 500 INTERNAL SERVER ERROR if one occurs

###GET api.polly.com/user/polls###
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Returns 200 OK for success, containing a list of polls and update time in JSON in the response body:
```
#!json
{
	"polls": [
		{
			"poll_id": 1123,
			"last_updated": 1424980524
		},

        ...
	]
}
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.
* Returns 500 INTERNAL SERVER ERROR if one occurs

###GET api.polly.com/poll:###
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Reads the list of poll identifiers from the GET parameter [id]. Example: GET http://api.polly.com/poll?id=0&id=1&id=2
* Accepts a maximum of 10 poll identifiers.
* Returns 200 OK for success, containing a list of polls and update time in JSON in the response body:
```
#!json
    {
        "polls" : [
            {
            "meta_data" : {
                "id" : 283
                "creation_date" : 1073029382
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

        },

        ...]
    }
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.
* Returns 500 INTERNAL SERVER ERROR if one occurs