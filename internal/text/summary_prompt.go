package text

const SUMMARY_PROMPT = `
All cycles of the thought process have been finished.
Right now, you job is to provide the system with a complete response. The response must contain ONLY the information that you with the user to see. No extra thincking, just a summary of everything you came to! It should be short and most importantly **it must DIRECTLY ANSWER THE USERS ORIGINAL PROMPT**
Your answer will later be combined with those of other agents to provide the user with a final answer!

You must provide some form of answer! Only fill in the yaml response field. *BUT RESPOND IN YAML*!
`

