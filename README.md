**POST api.polly.com/api/v1/register:**

* Request body:
    * [phone_number] a valid (dutch) phone number
* Returns 200 OK for success
    * Returns 400 BAD REQUEST when not providing a valid phone number
    * Returns 500 INTERNAL SERVER ERROR if one occurs

**POST api.polly.com/api/v1/register/verify:**
* Request body:
    * [phone_number] a valid (dutch) phone number
    * [verification_token] the verification token (for always 'VERIFY')
    * [display_name] the user's display name
    * [device_type] 0 for android, 1 for iOS
    * [device_guid] the device guid for push messages, ignored for now
* Returns 200 OK for success, containing the user in JSON in the response body:
```json
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

**POST api.polly.com/api/v1/poll:**
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Request body should contain JSON of the form:
```json
        {
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
```json        
        {
            "meta_data": {
                "poll_id": 5,                                           <-- Server-side poll ID
                "creator_id": 1,                                        <-- Creator Id, thus the sender's user ID
                "creation_date": 1428946988,                            <-- Unix time
                "last_updated": 1428946988,                             <-- Unix time
                "title": "Filmpje doen"
            },
            "question": {
                "id": 5,                                                <-- Server-side question ID
                "type": 0,
                "title": "Naar welke film gaan we vanavond?"
            },
            "options": [
                {
                    "id": 29,                                            <-- Server-side option ID
                    "question_id": 5,
                    "value": "The imitation game"
                },
                {
                    "id": 30,                                            <-- Server-side option ID
                    "question_id": 5,
                    "value": "American Sniper"
                },
                {
                    "id": 31,                                            <-- Server-side option ID
                    "question_id": 5,
                    "value": "Jupiter Ascend"
                }
            ],
            "votes": [],                                                 <-- Always empty
            "participants": [
                {
                    "id": 1,                                 
                    "phone_number": "0622197479",                        <-- User's phone number filled in
                    "display_name": "Bryan Eikema"                       <-- User's display name filled in
                }
            ]
        }
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.
* Returns 401 UNAUTHORIZED if no authentication is provided.
* Returns 500 INTERNAL SERVER ERROR if one occurs

**GET api.polly.com/api/v1/poll/xx:**
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Replace xx with the server-side identifier of the poll you're requesting.
* Returns 200 OK for success, containing the poll in JSON in the response body:
```json
        {
            "meta_data": {
                "poll_id": 1,
                "creator_id": 1,
                "creation_date": 1428938112,
                "last_updated": 1428938164,
                "title": "Filmpje doen"
            },
            "question": {
                "id": 1,
                "type": 0,
                "title": "Naar welke film gaan we vanavond?"
            },
            "options": [
                {
                    "id": 1,
                    "question_id": 1,
                    "value": "The imitation game"
                },
                {
                    "id": 2,
                    "question_id": 1,
                    "value": "American Sniper"
                },
                {
                    "id": 3,
                    "question_id": 1,
                    "value": "Jupiter Ascend"
                }
            ],
            "votes": [
                {
                    "id": 5,                                             <-- Server-side ID
                    "option_id": 2,                                      <-- Corresponding server-side option ID 
                    "user_id": 1,                                        <-- User-ID of the user who voted
                    "creation_date": 1428938164                          <-- Unix time
                }
            ],
            "participants": [
                {
                    "id": 1,
                    "phone_number": "0622197479",
                    "display_name": "Bryan Eikema"
                }
            ]
        }   
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.
* Returns 401 UNAUTHORIZED if no authentication is provided.
* Returns 403 FORBIDDEN when trying to access a poll in which the authorized user is no participant
* Returns 500 INTERNAL SERVER ERROR if one occurs

**GET api.polly.com/api/v1/user/polls**
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Accepts a page number as a GET parameter. Example: GET http://api.polly.com/user/polls?page=2
* Returns 200 OK for success, containing a list of polls and update time in JSON in the response body:
```json
        {
            "polls": [
                {
                    "poll_id": 1,
                    "last_updated": 1428938164
                },
                 
                ...
            ],
            "page": 1,                                                   <-- The page number of this response
            "page_size": 20,                                             <-- The maximum page size
            "num_results": 3,                                            <-- The number of results returned in this request
            "total": 3                                                   <-- The total number of results on the server
        }
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.
* Returns 401 UNAUTHORIZED if no authentication is provided.
* Returns 500 INTERNAL SERVER ERROR if one occurs

**GET api.polly.com/api/v1/poll:**
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Reads the list of poll identifiers from the GET parameter [id]. Example: GET http://api.polly.com/poll?id=0&id=1&id=2
* Accepts a maximum number of identifiers, more than the maximum will result in a 400 BAD REQUEST, this maximum is equal to the page size of GET /user/polls.
* Returns 200 OK for success, containing a list of polls and update time in JSON in the response body:
```json
        {
            "polls" : [
                {
                    "meta_data": {
                        "poll_id": 1,
                        "creator_id": 1,
                        "creation_date": 1428938112,
                        "last_updated": 1428938164,
                        "title": "Filmpje doen"
                    },
                    "question": {
                        "id": 1,
                        "type": 0,
                        "title": "Naar welke film gaan we vanavond?"
                    },
                    "options": [
                        {
                            "id": 1,
                            "question_id": 1,
                            "value": "The imitation game"
                        },
                        {
                            "id": 2,
                            "question_id": 1,
                            "value": "American Sniper"
                        },
                        {
                            "id": 3,
                            "question_id": 1,
                            "value": "Jupiter Ascend"
                        }
                    ],
                    "votes": [
                        {
                            "id": 5,                                       
                            "option_id": 2,                                    
                            "user_id": 1,                                  
                            "creation_date": 1428938164                       
                        }
                    ],
                    "participants": [
                        {
                            "id": 1,
                            "phone_number": "0622197479",
                            "display_name": "Bryan Eikema"
                        }
                    ]
                },   

            ...]
        }
```
* Returns 400 BAD REQUEST if information is wrong, incomplete or absent.
* Returns 401 UNAUTHORIZED if no authentication is provided.
* Returns 403 FORBIDDEN when trying to access a poll in which the authorized user is no participant
* Returns 500 INTERNAL SERVER ERROR if one occurs

**POST api.polly.com/api/v1/vote:**
* Requires the use of BasicAuth using the phone number and token as username and password, both preemptive and non-preemptive are supported.
* Request body should contain JSON of the form:
```json
        {
            "type" : 0,                                                  <-- 0 = new option, 1 = upvote
            "id" : 10,                                                   <-- Contains question ID for a new option, option ID for an upvote
            "value" : "New option"                                       <-- Contains the value of the new option, can be omitted for an upvote
        }

```
* Returns 200 OK for success, containing inserted vote and, if appropriate, the inserted option:
```json
        {
            "option": {                                                   <-- Omitted when type was 0
                "id": 53,                                                 <-- Server-side ID of the new option
                "question_id": 1,
                "value": "New option"
            },
            "vote": {
                "id": 1352,                                               <-- Server-side ID of the vote
                "option_id": 53,
                "user_id": 1,                                             <-- Server-side ID of the voter
                "creation_date": 1428949613                               <-- Unix time
            }
        }
```
* Returns 401 UNAUTHORIZED if no authentication is provided.
* Returns 403 FORBIDDEN when trying to access a poll in which the authorized user is no participant
* Returns 500 INTERNAL SERVER ERROR if one occurs



