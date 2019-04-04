import java.io.BufferedReader;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;

public class JavaTest {

    public static void main(String[] args) throws IOException {
        InputStream in = new FileInputStream("data");
        try (BufferedReader reader = new BufferedReader(new InputStreamReader(in))) {
            String[] raw = reader.readLine().split(",");
            int[] data = parseIntArray(raw);

            long start = System.currentTimeMillis();
            
            quickSort(data);

            System.out.print(System.currentTimeMillis() - start);
            System.out.print("ms ");
            System.out.println(check(data));
        }
    }

    static int[] parseIntArray(String[] raw) {
        int[] ret = new int[raw.length];
        for (int i = 0, limit = raw.length; i < limit; i++) {
            ret[i] = Integer.parseInt(raw[i], 10);
        }
        return ret;
    }

    static void quickSort(int[] data) {
        quickSort(data, 0, data.length - 1);
    }

    static void quickSort(int[] data, int start, int end) {
        // System.out.println(Arrays.toString(data) + "-" + start + ":" + end);
        if (end == start + 1) {
            if (data[start] > data[end]) {
                swap(data, start, end);
            }
        } else if (end > start) {
            int v = data[start];
            int i = start + 1, j = end;
            for (; i < j; i++) {
                if (data[i] > v) {
                    for (; i < j; j--) {
                        if (data[j] < v) {
                            swap(data, i, j);
                            break;
                        }
                    }
                    if (i == j) {
                        if (i > start + 1) {
                            swap(data, start, i - 1);
                            quickSort(data, start, i - 1);
                        }
                        quickSort(data, i, end);
                        return;
                    }
                }
            }
            if (data[i] >= v) {
                i--;
            }
            swap(data, start, i);
            quickSort(data, start, i);
            quickSort(data, i + 1, end);
        }
    }

    static void swap(int[]data, int a, int b) {
        data[a] = data[a] ^ data[b];
        data[b] = data[a] ^ data[b];
        data[a] = data[a] ^ data[b];
    }

    static boolean check(int[] data) {
        for (int i = 0, limit = data.length - 1; i < limit; i++) {
            if (data[i] > data[i + 1])
                return false;
        }
        return true;
    }

    

    public void testString() {
        String s = "asd1234551125";
        StringBuilder builder = new StringBuilder(s.length());
        builder.append(s);
        builder.reverse();
        s = builder.toString();

        char[] tmp = new char[s.length()];
        for (int i = 0, j = s.length() - 1; i <= j; i++, j--) {
            tmp[j] = s.charAt(i);
            tmp[i] = s.charAt(j);
        }
        String reversed = new String(tmp);
        System.out.println(reversed);
    }

    private final Integer lock = 0;

    public void testLock() throws InterruptedException {
        Thread t = new Thread(() -> ok());
        t.start();
        waiting();
    }

    private void waiting() throws InterruptedException {
        synchronized (lock) {
            System.out.println("waiting");
            lock.wait();
        }
    }

    private void ok() {
        try {
            Thread.sleep(2000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        System.out.println("notifying");
        synchronized(lock) {
            lock.notify();
        }
    }

}
