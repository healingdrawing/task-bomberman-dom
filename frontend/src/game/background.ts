// Trying to make random background position 
// TODO: How to implement this to #game_object ? 
function getRandomBackgroundPosition(): string {
    const positions: string[] = ['top', 'bottom', 'left', 'right', 'center'];
    const randomIndex: number = Math.floor(Math.random() * positions.length);
    return positions[randomIndex];
  }
  
  document.addEventListener("DOMContentLoaded", () => {
    const gameObjects: HTMLElement | null = document.getElementById('game_objects');
  
    if (gameObjects) {
      // Apply the random background position to the element
      gameObjects.style.backgroundPosition = getRandomBackgroundPosition();
    }
  });
  