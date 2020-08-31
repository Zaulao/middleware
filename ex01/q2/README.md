Questão 2: Um número indefinido (máximo de 10) de threads (consumidores) podem
executar a função int consumir(), e apenas um thread (produtor) pode chamar a
função void produzir (int n). Os comportamentos destas duas funções são
descritos a seguir:
consumir() é bloqueante enquanto espera que o produtor invoque a função
produzir(). A função produzir(), quando chamada, desbloqueia todos os
consumidores bloqueados por consumir(). A função consumir() deve retornar
aos consumidores que estiverem bloqueados o valor passado como argumento para a
função produzir(n).
A função produzir(n), no caso de haver consumidores bloqueados por consumir(),
deve desbloqueá-los e retornar para os consumidores o valor n. Caso não haja
nenhum consumidor bloqueado, o produtor deve ser bloqueado até que um
consumidor execute a função consumir().
Exemplo 1 (possível sequência de invocações):
 
C1 (consumidor) invoca consumir() e fica bloqueado;
C2 (consumidor) invoca consumir() e fica bloqueado;
P (produtor) invoca produzir(13) e isto desbloqueia os consumidores C1 e C2.
Neste caso, a função consumir() retorna o valor 13 para C1 e C2.
 
Exemplo 2 (possível sequência de invocações):
 
P (produtor) invoca produzir(23) e bloqueia;
C1 (consumidor) invoca consumir() que desbloqueia P e retorna valor 23 para C1;
C2 invoca consumir() e bloqueia.

