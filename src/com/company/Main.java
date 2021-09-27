package com.company;

import java.sql.Array;
import java.util.ArrayList;
import java.util.BitSet;
import java.util.Scanner;
import java.util.Arrays;

public class Main {
    public static String[] valid = new String[]{"+","-","*","/","1","2","3","4","5","6","7","8","9"};

    public static void main(String[] args) {
        System.out.println("Welcome to the Rule 110 Regular Expression Calculator!");
        System.out.println("         ------------------------------------         ");
        System.out.println();
        Scanner input = new Scanner(System.in);

        boolean check = false;
        String p = "x";
        while(!check) { //check is false
            System.out.println("NOTE: expressions can only have 1 operator between terms!");
            System.out.print("Please enter an expression with only numbers and (+),(-),(*),(/): ");
            ArrayList<String> exp = new ArrayList<String>(Arrays.asList(input.nextLine().split("")));
            exp.removeIf(x -> (x.equals(" ")));

            if (Arrays.asList(Arrays.copyOfRange(valid, 0, 4)).contains(exp.get(0)) || Arrays.asList(Arrays.copyOfRange(valid, 0, 4)).contains(exp.get(exp.size()-1))) {
                System.out.println(); continue; }

            for (String x : exp) {
                for (String y : valid) {
                    if(x.equals(y)) {check = true; break;}
                    check = false;
                }
                if(!check) {break;}
                if(Arrays.asList(Arrays.copyOfRange(valid, 0,4)).contains(p) && Arrays.asList(Arrays.copyOfRange(valid, 0,4)).contains(x)) {
                    check=false; break;}
                p=x;
            }
            System.out.println();
        }
        System.out.println("System Works!");
        //TODO add check for no operator
    }

    public void rule110Sim(BitSet initial, int iter){
        BitSet[] matrix = new BitSet[iter];
        matrix[0] = initial;
        for(int i = 1; i < iter; i++) {

        }
    }

    public String[] simplify(ArrayList<String> exp) {
        ArrayList<String> temp = new ArrayList<String>(exp);
        for(int i = 0; i < exp.size(); i++) {
            if(exp.get(i).equals("*") || exp.get(i).equals("/")){

            }
        }
        return null;
        // recursion to simplify
        // this method will only make the problem into addition and subtraction
    }

}
