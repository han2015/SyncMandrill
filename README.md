# SyncMandrill
 syncs templats from mailchimp account to our test/dev mandrill accounts.
 Production templates could updated by providing two environment variables PROD_KEY.
 
### install
 1: cd SyncMandrill
 
 2: go install 
 
### how to use it.
  1: Create a config.json for your project.
  
  2: exec SyncMandrill .


### config.json
 `{
	"MailChimp":{
		"APIKey":"111111111111"
	},
	"Official":{"Name":"official account on mandrill ","APIKey":"222222222"},
	"Accounts":[
		{"Name":"test account on mandrill ","APIKey":"222222222"}
	],
	"Slugs":{"exampleHan":"exampleHan-en"}
}`