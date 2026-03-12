from tempura import batter


@batter
def main(hi: str = "asd"):
    print("Hello from tempura!")


if __name__ == "__main__":
    main()
