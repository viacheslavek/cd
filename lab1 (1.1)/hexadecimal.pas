program HexExample;

const
  HexNumber = $A3;  { Это шестнадцатеричное число, равное 163 в десятичной системе }

var
  DecimalNumber: Integer;  { Десятичное число для сложения }
  Sum: Integer;  { Переменная для хранения суммы }

begin
  DecimalNumber := 10;  { Инициализация десятичного числа }
  Sum := HexNumber + DecimalNumber;  { Сложение шестнадцатеричного и десятичного числа }
  writeln('Сумма ', HexNumber, ' (в шестнадцатеричной) и ', DecimalNumber, ' (в десятичной) равна ', Sum);
end.
