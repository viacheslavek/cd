% Лабораторная работа № 3 «Обобщённые классы в Scala»
% 21 марта 2024 г.
% Вячеслав Локшин, ИУ9-61Б

# Цель работы
Целью данной работы является приобретение навыков разработки обобщённых классов на языке Scala
с использованием неявных преобразований типов.


# Индивидуальный вариант
Класс ArithProgression[T], представляющий арифметическую прогрессию 
с элементами типа T и операцией вычисления i-го члена и суммы первых n членов.
Арифметическая прогрессия должны быть задана начальным членом и разностью.
В качестве типа T может выступать числовой тип или строка. В случае строки сложение означает конкатенацию.

# Реализация

```scala
import scala.annotation.unused
import scala.collection.mutable

trait ArithProgression[T] {
  def get_i_member(n: Int, start: T, difference: T): T
  def get_sum_i_members(n: Int, start: T, difference: T): T
}

@unused // баг IDE
object ArithProgression {
  implicit object StringArithProgression extends ArithProgression[String] {
    def get_i_member(n: Int, start: String, difference: String): String = {
      start + difference * n
    }

    def get_sum_i_members(n: Int, start: String, difference: String): String = {
      val sb = new mutable.StringBuilder
      val curSb = new mutable.StringBuilder(start)

      for (_ <- 0 until n) {
        sb.append(curSb)
        sb.append(difference)
        curSb.append(difference)
      }
      sb.toString

    }
  }

  implicit def fractionalArithProgression[T: Numeric]: ArithProgression[T] = new ArithProgression[T] {
    def get_i_member(n: Int, start: T, difference: T): T = {
      Numeric[T].plus(
        start, implicitly[Numeric[T]].times(difference, implicitly[Numeric[T]].fromInt(n - 1))
      )
    }

    def get_sum_i_members(n: Int, start: T, difference: T): T = {
      // INFO: (a1+an)*n/2

      if (n % 2 == 0) {
        val newN = n / 2
        implicitly[Numeric[T]].times(
          implicitly[Numeric[T]].fromInt(newN),
          Numeric[T].plus(start, get_i_member(n, start, difference)),
        )
      } else {
        val newN = (n - 1) / 2
        Numeric[T].plus(
          implicitly[Numeric[T]].times(
            implicitly[Numeric[T]].fromInt(newN),
            Numeric[T].plus(start, get_i_member(n-1, start, difference))),
          get_i_member(n, start, difference)
        )
      }
    }
  }


}

class ArithProgressionObj[T](val start: T, val difference: T)(implicit ev: ArithProgression[T]) {
  def Get_i_member(n: Int): T = {
    println("n: " + n + " | start: " + start + " | diff: " + difference)
    ev.get_i_member(n, start, difference)
  }

  def Get_sum_i_members(n: Int): T = {
    println("n: " + n + " | start: " + start + " | diff: " + difference)
    ev.get_sum_i_members(n, start, difference)
  }
}

object Main extends App {

  def DoWork[T](start: T, difference: T, n: Int)(implicit ev: ArithProgression[T]): Unit = {
    val arith = new ArithProgressionObj[T](start, difference)
    println("n-тый член арифметической прогрессии: "+arith.Get_i_member(n).toString)
    println("Сумма первых n членов арифметической прогрессии: "+arith.Get_sum_i_members(n).toString)

  }

  println("START\n")

  DoWork[Int](0, 1, 10)
  DoWork[Int](0, 1, 9)

  DoWork[Double](0.5, 10.1, 5)

  DoWork[String]("A", "B", 5)

  println("\nFINISH")

}

```

# Тестирование

Результат запуска программы:

```
START

n: 10 | start: 0 | diff: 1
n-тый член арифметической прогрессии: 9
n: 10 | start: 0 | diff: 1
Сумма первых n членов арифметической прогрессии: 45
n: 9 | start: 0 | diff: 1
n-тый член арифметической прогрессии: 8
n: 9 | start: 0 | diff: 1
Сумма первых n членов арифметической прогрессии: 36
n: 5 | start: 0.5 | diff: 10.1
n-тый член арифметической прогрессии: 40.9
n: 5 | start: 0.5 | diff: 10.1
Сумма первых n членов арифметической прогрессии: 103.5
n: 5 | start: A | diff: B
n-тый член арифметической прогрессии: ABBBBB
n: 5 | start: A | diff: B
Сумма первых n членов арифметической прогрессии: ABABBABBBABBBBABBBBB

FINISH
```

# Вывод
В ходе данной лабораторной работы был получен опыт разработки обобщённых классов на языке Scala
с использованием неявных преобразований типов.
Интересно, какие большие возможности предоставляет язык Scala.
