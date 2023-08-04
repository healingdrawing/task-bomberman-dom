class ScreenPrepare {

  cspx = 96 // cell size in px
  //todo: potential way to improve performance
  /*
  hope this will be not needed, because the field includes only 49 cells maximum 7x7. it is not 1280x720 elements, when you recalculate every pixel on screen separately.
  
  Instead of maps we can use arrays,
  but it will be less readable(empty elements will present inside),
  and extra stuff can be required, to make indexes from x and y
  (multiply by 10 at least x coordinate, before generate the index x*10+y), 
  and perhaps, convert cells number, but generally the one-dimensional array approach is faster
  */
  weak_obstacles = new Map<string, HTMLDivElement>() // keep weak obstacles, to manipulate them later
  power_ups = new Map<string, HTMLDivElement>() // keep powerups, to manipulate them later
  bombs = new Map<string, HTMLDivElement>() // keep bombs, to manipulate them later
  players = new Map<string, HTMLDivElement>() // keep players, to manipulate them later
  explosions = new Map<string, HTMLDivElement>() // keep explosions, to manipulate them later

  /** prebuild_game_field() generates all required objects to display the game
   * and place them on screen, so later we can manipulate them, just show/hide/move,
   * without creating new objects, so we can save memory and cpu time,
   * and browser do not need even rebuild the DOM tree,
   * just change some css properties, generally styles of the divs
   */
  async prebuild_game_field(): Promise<void> {
    console.log("prebuild_game_field")
    // build strong obstacles, can not be destroyed
    const game_strong_obstacles = document.getElementById("game_strong_obstacles") as HTMLDivElement;
    const dx = [1, 3, 5, 1, 3, 5, 1, 3, 5]
    const dy = [1, 1, 1, 3, 3, 3, 5, 5, 5]

    for (let i = 0; i < 9; i++) {
      const div = document.createElement("div");
      div.classList.add("cell"); // for static objects
      div.classList.add("strong_obstacle");
      div.style.left = `${dx[i] * this.cspx}px`;
      div.style.top = `${dy[i] * this.cspx}px`;
      game_strong_obstacles.appendChild(div);
    }

    // build weak obstacles, can be destroyed. CCW from right center
    const game_weak_obstacles = document.getElementById("game_weak_obstacles") as HTMLDivElement;
    dx.length = 0
    dy.length = 0
    dx.push(6, 3, 0, 3, 4, 3, 2, 3) // first 4 are for external, last 4 for internal
    dy.push(3, 0, 3, 6, 3, 2, 3, 4)

    for (let i = 0; i < 8; i++) {
      const div = document.createElement("div");
      div.classList.add("cell"); // for static objects
      div.classList.add("weak_obstacle"); // todo: remove this later
      div.style.left = `${dx[i] * this.cspx}px`;
      div.style.top = `${dy[i] * this.cspx}px`;
      game_weak_obstacles.appendChild(div);
      this.weak_obstacles.set(`${dx[i]}${dy[i]}`, div) // always less then 10, so no need space between x and y
    }

    // build power_ups, they can be placed under weak obstacles, and will be shown after weak obstacle is destroyed, and server will send us info about it

    const game_power_ups = document.getElementById("game_power_ups") as HTMLDivElement;
    for (let i = 0; i < 8; i++) {
      const div = document.createElement("div");
      div.classList.add("absolute"); // for animated objects, to do not lock size
      div.classList.add("power_up1"); //todo: remove this later
      div.style.left = `${dx[i] * this.cspx}px`;
      div.style.top = `${dy[i] * this.cspx}px`;
      game_power_ups.appendChild(div);
      this.power_ups.set(`${dx[i]}${dy[i]}`, div) // always less then 10, so no need space between x and y
    }

    // build bombs
    const game_bombs = document.getElementById("game_bombs") as HTMLDivElement;
    dx.length = 0
    dy.length = 0
    dx.push(
      0, 1, 2, 3, 4, 5, 6,
      0, 2, 4, 6,
      0, 1, 2, 3, 4, 5, 6,
      0, 2, 4, 6,
      0, 1, 2, 3, 4, 5, 6,
      0, 2, 4, 6,
      0, 1, 2, 3, 4, 5, 6,
    )
    dy.push(
      0, 0, 0, 0, 0, 0, 0,
      1, 1, 1, 1,
      2, 2, 2, 2, 2, 2, 2,
      3, 3, 3, 3,
      4, 4, 4, 4, 4, 4, 4,
      5, 5, 5, 5,
      6, 6, 6, 6, 6, 6, 6,
    )

    for (let i = 0; i < 40; i++) {
      const div = document.createElement("div");
      div.classList.add("absolute"); // for animated objects, to do not lock size
      // div.classList.add("player1_bomb"); //todo: remove this later
      div.style.left = `${dx[i] * this.cspx}px`;
      div.style.top = `${dy[i] * this.cspx}px`;
      game_bombs.appendChild(div);
      this.bombs.set(`${dx[i]}${dy[i]}`, div) // always less then 10, so no need space between x and y
    }

    // build explosions, coordinates the same as bombs
    const game_explosions = document.getElementById("game_explosions") as HTMLDivElement;

    for (let i = 0; i < 40; i++) {
      const div = document.createElement("div");
      div.classList.add("absolute"); // for animated objects, to do not lock size
      // div.classList.add("explosion"); //todo: remove this later
      div.style.left = `${dx[i] * this.cspx}px`;
      div.style.top = `${dy[i] * this.cspx}px`;
      game_explosions.appendChild(div);
      this.explosions.set(`${dx[i]}${dy[i]}`, div) // always less then 10, so no need space between x and y
    }

    // build players
    const game_players = document.getElementById("game_players") as HTMLDivElement;
    dx.length = 0
    dy.length = 0
    dx.push(0, 6, 0, 6)
    dy.push(0, 6, 6, 0)

    for (let i = 0; i < 4; i++) {
      const div = document.createElement("div");
      div.classList.add("absolute"); // for animated objects, to do not lock size
      div.classList.add(`player${i + 1}_move`); //todo: remove this later
      div.style.left = `${dx[i] * this.cspx}px`;
      div.style.top = `${dy[i] * this.cspx}px`;
      game_players.appendChild(div);
      this.players.set(`${dx[i]}${dy[i]}`, div) // always less then 10, so no need space between x and y
    }
  }

}

export const screen_prepare = new ScreenPrepare()