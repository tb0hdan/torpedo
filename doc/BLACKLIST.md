# [Back to main doc](../README.md)

# Sender/Message blacklist

## Create collection

```
db.createCollection("blackListItems")
```


## Blacklisting by message

```
db.blackListItems.insertOne({type: 'message', pattern: 'cheap advertising', matches: NumberLong(0)})
```


## Blacklisting by sender


```
db.blackListItems.insertOne({type: 'sender', pattern: 'nastyspammer@jabber.org', matches: NumberLong(0)})
```
