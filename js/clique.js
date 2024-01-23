document.addEventListener('keydown', (event) => {
    let audio;
    switch (event.key) {
      case "Enter":
        audio = new Audio('/audio/enter.mp3');
        break;
      case "Backspace":
        audio = new Audio('/audio/backspace.mp3');
        break;
      default:
        audio = new Audio('/audio/key.mp3');
    }
    audio.play();
  });