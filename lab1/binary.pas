program BinExample;

const
  BinNumber = 0b10100011;  { Это двоичное число, равное 163 в десятичной системе }

var
  DecimalNumber: Integer;
  Sum: Integer;

begin
  DecimalNumber := 10;
  Sum := BinNumber + DecimalNumber;
  writeln('Сумма ', BinNumber, ' (в двоичной) и ', DecimalNumber, ' (в десятичной) равна ', Sum);
end.
