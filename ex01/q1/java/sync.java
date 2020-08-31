import java.util.concurrent.*;
import java.util.concurrent.TimeUnit;

public class sync {
	public static void main(String args[]) {
		final ponte_sync ponte = new ponte_sync();

		Thread indo = new Thread( new Runnable() {   
			public void run() {

				while(true)
				{
					try
					{
						TimeUnit.SECONDS.sleep((long)(Math.random()*10));
						System.out.println("Carro indo está tentando cruzar a ponte");
						ponte.atravessar("Carro indo");
					}
					catch(InterruptedException iex)
					{
						iex.printStackTrace();
					}					
				}
			}
		});
		Thread voltando = new Thread( new Runnable() {   
			public void run() {

				while(true)
				{
					try
					{
						TimeUnit.SECONDS.sleep((long)(Math.random()*10));
						System.out.println("Carro voltando está tentando cruzar a ponte");
						ponte.atravessar("Carro voltando");
					}
					catch(InterruptedException iex)
					{
						iex.printStackTrace();
					}					
				}
			}
		});
		indo.start();
		voltando.start();
	}

}

class ponte_sync{
	private int atravessando;

	public ponte_sync() {
		atravessando = 0;
	}

	public synchronized void atravessar(String nome) throws InterruptedException{
		System.out.println(nome + " está cruzando a ponte.");
		this.atravessando++;
		TimeUnit.SECONDS.sleep(3);
		if(atravessando > 1) {
			System.out.println("ACIDENTE!!!");
			System.out.println(nome + " não conseguiu atravessar");
		}else {
		this.atravessando--;
		System.out.println(nome + " atravessou");
		}
	}
}