# include <iostream>

using namespace std;

double square(double x) {
    return x * x;
}

void print_square(double x) {
    cout << "the square of " << x << " is " << square(x) << "\n";
}

int main() {
    cout << "Hello, Wayne's World! \n";
    print_square(1.234);
    return 0;
}