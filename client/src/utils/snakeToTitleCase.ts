export default function snakeToTitleCase(snake: string): string {
    return snake
        .split("_")
        .filter((sub) => sub != "")
        .map((sub) => sub.at(0)?.toUpperCase() + sub.slice(1))
        .join(" ")
}
