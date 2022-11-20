# IS MY PASSWORD SAFE

## What is it

Return how many times a password has been leaked. It uses haveibeenpwned API.
It uses https://api.pwnedpasswords.com/range/ to get a list of hashes that has the first 5 bytes equal to the first 5 bytes of the typed password, then iterates over the list to find the correspondent hash.

## Input

HTTP POST

```
{
    "password": "Password1"
}
```

## Output

```
{
    "password": "Password1",
    "occurrences": "118930",
    "sha1": "70CCD9007338D6D81DD3B6271621B9CF9A97EA00"
}
```
