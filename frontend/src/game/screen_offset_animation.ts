
export interface OffsetAnimationParameters {
  div: HTMLElement;
  player_number: number;
  fromX: number;
  toX: number;
  fromY: number;
  toY: number;
  cells_per_second: number;
  cspx: number;
}

export function offset_div_animation({
  div,
  player_number,
  fromX,
  toX,
  fromY,
  toY,
  cells_per_second,
  cspx,
}: {
  div: HTMLElement;
  player_number: number;
  fromX: number;
  toX: number;
  fromY: number;
  toY: number;
  cells_per_second: number;
  cspx: number;
}) {
  const distanceX = toX - fromX;
  const distanceY = toY - fromY;
  const sum_offset = Math.sqrt(distanceX * distanceX + distanceY * distanceY) * cspx;

  const ani_time_ms = sum_offset / (cells_per_second * cspx) * 1000;

  const intervalDuration = 16;
  const startTime = Date.now()

  const animationInterval = setInterval(() => {

    const progress = (Date.now() - startTime) / ani_time_ms;

    if (progress >= 1) {
      clearInterval(animationInterval);
      div.classList.remove(`player${player_number}_move`)
      div.style.left = `${toX * cspx}px`;
      div.style.top = `${toY * cspx}px`;
      return;
    }

    const currentX = (fromX + distanceX * progress) * cspx;
    const currentY = (fromY + distanceY * progress) * cspx;

    div.style.left = `${currentX}px`;
    div.style.top = `${currentY}px`;
    div.classList.add(`player${player_number}_move`)
  }, intervalDuration);
}

/** returns cell_x of div, but not integer */
export function get_from_x(div: HTMLDivElement, cspx: number): number {
  const from_x = parseInt(div.style.left) / cspx
  return from_x
}

/** returns cell_y of div, but not integer */
export function get_from_y(div: HTMLDivElement, cspx: number): number {
  const from_y = parseInt(div.style.top) / cspx
  return from_y
}