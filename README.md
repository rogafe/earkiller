tranform to dca 

```batch
ffmpeg -i file.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > file.dca
```