package hanzofs.examples;

import hanzofs.client.FilerClient;
import hanzofs.client.FilerInputStream;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

public class ExampleReadFile {

    public static void main(String[] args) throws IOException {

        FilerClient filerClient = new FilerClient("localhost", 18888);

        long startTime = System.currentTimeMillis();
        parseZip("/Users/chris/tmp/test.zip");

        long startTime2 = System.currentTimeMillis();

        long localProcessTime = startTime2 - startTime;

        FilerInputStream filerInputStream = new FilerInputStream(
                filerClient, "/test.zip");
        parseZip(filerInputStream);

        long swProcessTime = System.currentTimeMillis() - startTime2;

        System.out.println("Local time: " + localProcessTime);
        System.out.println("Hanzo S3 time: " + swProcessTime);

    }

    public static void parseZip(String filename) throws IOException {
        FileInputStream fileInputStream = new FileInputStream(filename);
        parseZip(fileInputStream);
    }

    public static void parseZip(InputStream is) throws IOException {
        ZipInputStream zin = new ZipInputStream(is);
        ZipEntry ze;
        while ((ze = zin.getNextEntry()) != null) {
            System.out.println(ze.getName());
        }
    }
}
