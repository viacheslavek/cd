% Лабораторная работа № 1.1. Раскрутка самоприменимого компилятора
% 9 февраля 2024 г.
% Локшин Вячеслав, ИУ9-61Б

# Цель работы
Целью данной работы является ознакомление с раскруткой самоприменимых компиляторов
на примере модельного компилятора.

# Индивидуальный вариант
Компилятор BeRo. Добавить в язык двоичные константы вида 0b10010.

# Реализация

Различие между файлами `pcom.pas` и `pcom2.pas`:

```diff
@@ -594,10 +594,24 @@
 begin
  Num:=0;
  if ('0'<=CurrentChar) and (CurrentChar<='9') then begin
-  while ('0'<=CurrentChar) and (CurrentChar<='9') do begin
-   Num:=(Num*10)+(ord(CurrentChar)-ord('0'));
-   ReadChar;
-  end;
+
+   // Проверка на двоичное число
+   if CurrentChar = '0' then begin
+     ReadChar;
+     if (CurrentChar = 'b') then begin
+       ReadChar;
+         while (CurrentChar = '0') or (CurrentChar = '1') do begin
+           Num := (Num * 0b10) + (ord(CurrentChar) - ord('0'));
+           ReadChar;
+         end;
+     end;
+   end else begin
+     Num := 0;
+     while ('0'<=CurrentChar) and (CurrentChar<='9') do begin
+       Num:=(Num*10)+(ord(CurrentChar)-ord('0'));
+       ReadChar;
+     end;
+   end;
  end else if CurrentChar='$' then begin
   ReadChar;
   while (('0'<=CurrentChar) and (CurrentChar<='9')) or
@@ -652,6 +666,7 @@
    s:=s+1;
   end;
  end else if (('0'<=CurrentChar) and (CurrentChar<='9')) or (CurrentChar='$') then begin
+  { INFO: здесь ничего менять не надо - уже зайду в распознование числа по 0 }
   CurrentSymbol:=TokNumber;
   CurrentNumber:=ReadNumber;
  end else if CurrentChar=':' then begin
```

Различие между файлами `pcom2.pas` и `pcom3.pas`:

```diff
@@ -168,7 +168,7 @@
       SymFUNC=49;
       SymPROC=50;
 
-      IdCONST=0;
+      IdCONST=0b0;
       IdVAR=1;
       IdFIELD=2;
       IdTYPE=3;
@@ -601,7 +601,7 @@
      if (CurrentChar = 'b') then begin
        ReadChar;
          while (CurrentChar = '0') or (CurrentChar = '1') do begin
-           Num := (Num * 0b10) + (ord(CurrentChar) - ord('0'));
+           Num := (Num * 2) + (ord(CurrentChar) - ord('0'));
            ReadChar;
          end;
      end;

```

# Тестирование

Тестовый пример:

```pascal
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

```

Вывод тестового примера на `stdout`

```
>$ ./bin

Сумма 163 (в двоичной) и 10 (в десятичной) равна 173

```

# Вывод
Потренировал навыки чтения и анализа кода и ознакомился с раскруткой самоприменимых компиляторов
на примере модельного компилятора.