package hanzofs.examples;

import hanzofs.client.FilerClient;
import hanzofs.client.FilerInputStream;
import hanzofs.client.FilerOutputStream;

import java.io.IOException;
import java.io.InputStream;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

public class ExampleWriteFile {

    public static void main(String[] args) throws IOException {

        FilerClient filerClient = new FilerClient("localhost", 18888);

        FilerInputStream filerInputStream = new FilerInputStream(filerClient, "/test.zip");
        unZipFiles(filerClient, filerInputStream);

    }

    public static void unZipFiles(FilerClient filerClient, InputStream is) throws IOException {
        ZipInputStream zin = new ZipInputStream(is);
        ZipEntry ze;
        while ((ze = zin.getNextEntry()) != null) {

            String filename = ze.getName();
            if (filename.indexOf("/") >= 0) {
                filename = filename.substring(filename.lastIndexOf("/") + 1);
            }
            if (filename.length()==0) {
                continue;
            }

            FilerOutputStream filerOutputStream = new FilerOutputStream(filerClient, "/test/"+filename);
            byte[] bytesIn = new byte[16 * 1024];
            int read = 0;
            while ((read = zin.read(bytesIn))!=-1) {
                filerOutputStream.write(bytesIn,0,read);
            }
            filerOutputStream.close();

            System.out.println(ze.getName());
        }
    }
}
