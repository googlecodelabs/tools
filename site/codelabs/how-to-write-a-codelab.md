summary: How to Write a Codelab
id: how-to-write-a-codelab
categories: Sample
tags: medium
status: Published 
authors: Zarin
Feedback Link: https://zarin.io

# How to Write a Codelab

<!-- ---------------------------------------------------------------------------- --> 
## Overview 
Duration: 1

### What You'll Learn 
- how to set the amount of time each slide will take to finish 
- how to include code snippets 
- how to hyperlink items 
- how to include images 
- other stuff

<!-- ---------------------------------------------------------------------------- --> 
## Setting Duration
Duration: 2

To indicate how long each slide will take to go through, set the `Duration` under each Heading 2 (i.e. `##`) to an integer. 
The integers refer to minutes. If you set `Duration: 4` then a particular slide will take 4 minutes to complete. 

The total time will automatically be calculated for you and will be displayed on the codelab once you create it. 

<!-- ---------------------------------------------------------------------------- --> 
## Code Snippets
Duration: 3

To include code snippets you can do a few things. 
- Inline highlighting can be done using the tiny tick mark on your keyboard: "`"
- Embedded code

### JavaScript

```javascript
{ 
  key1: "string", 
  key2: integer,
  key3: "string"
}
```

### Java

```java
for (statement 1; statement 2; statement 3) {
  // code block to be executed
}
```

<!-- ---------------------------------------------------------------------------- --> 
## Hyperlinking and Embedded Images

### Hyperlinking
[Youtube - Halsey Playlists](https://www.youtube.com/user/iamhalsey/playlists)

### Images
![alt-text-here](assets/puppy.jpg)

<!-- ---------------------------------------------------------------------------- --> 
## Other Stuff
Duration: 1

Checkout the official documentation here: [Codelab Formatting Guide](https://github.com/googlecodelabs/tools/blob/master/FORMAT-GUIDE.md)
