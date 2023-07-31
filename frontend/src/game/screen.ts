const cspx = 96 // cell size in px

export async function prebuild_game_field(): Promise<void> {
  console.log("prebuild_game_field")
  // build strong obstacles, can not be destroyed
  const game_strong_obstacles = document.getElementById("game_strong_obstacles") as HTMLDivElement;
  const dx = [1, 3, 5, 1, 3, 5, 1, 3, 5]
  const dy = [1, 1, 1, 3, 3, 3, 5, 5, 5]

  for (let i = 0; i < 9; i++) {
    const div = document.createElement("div");
    div.classList.add("cell");
    div.classList.add("strong_obstacle");
    div.style.left = `${dx[i] * cspx}px`;
    div.style.top = `${dy[i] * cspx}px`;
    game_strong_obstacles.appendChild(div);
  }
}