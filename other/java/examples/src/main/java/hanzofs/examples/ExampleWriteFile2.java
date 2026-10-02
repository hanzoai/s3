package hanzofs.examples;

import com.google.common.io.Files;
import hanzofs.client.FilerClient;
import hanzofs.client.FilerOutputStream;

import java.io.File;
import java.io.IOException;

public class ExampleWriteFile2 {

    public static void main(String[] args) throws IOException {

        FilerClient filerClient = new FilerClient("localhost", 18888);

        FilerOutputStream filerOutputStream = new FilerOutputStream(filerClient, "/test/1");
        Files.copy(new File("/etc/resolv.conf"), filerOutputStream);
        filerOutputStream.close();

    }

}
