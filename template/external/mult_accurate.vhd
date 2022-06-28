--Salim Ullah, Sanjeev Sripadraj Murthy, and Akash Kumar. 2018. SMApproxlib: library of FPGA-based approximate multipliers. In Proceedings of the 55th Annual Design Automation Conference (DAC '18). ACM, New York, NY, USA, Article 157, 6 pages. DOI: https://doi.org/10.1145/3195970.3196115

library IEEE;
use IEEE.STD_LOGIC_1164.ALL; 
use IEEE.NUMERIC_STD.ALL; 
use IEEE.STD_LOGIC_UNSIGNED.ALL;
library UNISIM;
use UNISIM.VComponents.all;

entity {{.EntityName}} is
generic (word_size: integer:={{.BitSize}}); 
Port ( 
A : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
B : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0));
end {{.EntityName}};
    
architecture Behavioral of {{.EntityName}} is

type data_type is array((word_size - 2 ) downto 0) of std_logic_vector((word_size-1) downto 0);
signal prop: data_type;
signal gen: data_type;
signal carries: data_type;
signal output: data_type;
type input_carries is array((word_size - 2 ) downto 0) of std_logic_vector(word_size/4 downto 0);
signal input_carry: input_carries;
type interm_results is array((word_size - 2) downto 0) of std_logic_vector(word_size downto 0);
signal chain: interm_results;
--signal prod : STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0);
--signal a : STD_LOGIC_VECTOR (word_size-1 downto 0);
--signal b : STD_LOGIC_VECTOR (word_size-1 downto 0);

begin
    
--    process (clk)
--    begin
--        if rising_edge(clk) then
--            a(7 downto 0) <= inp_a(7 downto 0);
--            b(7 downto 0) <= inp_b(7 downto 0);
--            product(15 downto 0) <= prod(15 downto 0);
--        end if;
--    end process;
    
    set_initial_carry:
    for car in 0 to (word_size - 2) generate
        input_carry(car)(0) <= '0';
    end generate set_initial_carry; 
    ----- LUTs Type A ------------
    Type_A: 
    for i in 0 to (word_size - 2) generate 
        lut_inst0: lut6_2 
        generic map(INIT => X"7888788880008000")
        port map(
        I0 => b(0), --h
        I1 => a(i+1), --c
        I2 => b(1), -- g
        I3 => a(i), -- d
        I4 => '1',
        I5 => '1',
        O5 => gen(0)(i),
        O6 => prop(0)(i)
        );
    end generate Type_A;
    
    lut_inst1: lut6_2 
    generic map(INIT => X"7888788880008000")
    port map(
    I0 => b(2), 
    I1 => a(word_size-2), 
    I2 => b(1), 
    I3 => a(word_size-1), 
    I4 => '1',
    I5 => '1',
    O5 => gen(0)(word_size-1) ,
    O6 => prop(0)(word_size-1)
    );
    ------------- End of Type A LUTs ---------
    carry_chain_type_A:										-- for word_size =8, this will loop will iterate two times (0 and 1)
    for j in 0 to (word_size/4 - 1) generate			--
        carry_inst0: CARRY4										--
        port map (													-- so first gen(0) and gen (1) are used which are 4 bits wide registers
            DI => gen(0)(j*4+3 downto j*4), 											--
            S => prop(0)(j*4+3 downto j*4),											--
            O => output(0)(j*4+3 downto j*4),										--
            CO => carries(0)(j*4+3 downto j*4),										-- output carries are stored in carries(0) and carries(1)
            CI => input_carry(0)(j),							-- first input carry should be input_carry(0)(0) and second should be input_carry(1)(0)
            CYINIT => '0'
            );
        input_carry(0)(j+1) <= carries(0)(j*4+3);
        chain(0)(j*4+3 downto j*4) <= output(0)(j*4+3 downto j*4);
    end generate carry_chain_type_A;
    chain(0)(word_size) <= carries(0)(word_size - 1);
    -----------------------------------------
    ----- LUTs Type B ------------
    Type_B:
    for t in 0 to (word_size - 4) generate
        B_steps:
            for u in 0 to (word_size - 3 - t) generate
                lut_inst2: lut6_2 
                    generic map(INIT => X"7878787880808080")
                    port map(
                    I0 => b(t + 2),
                    I1 => a(u),
                    I2 => chain(t)(u+1),
                    I3 => '1',
                    I4 => '1',
                    I5 => '1',
                    O5 => gen(t+1)(u), --- I am chainging the size of gen and prop from here. Gen and Prop are now 8 bit registers.
                    O6 => prop(t+1)(u)
                    );
        end generate B_steps;
        -- this is for the irregular LUT
        lut_inst3: lut6_2 
            generic map(INIT => X"7878787880808080")
            port map(
            I0 => b(t + 3),
            I1 => a(word_size-3-t),
            I2 => chain(t)(word_size-3-t+2),
            I3 => '1',
            I4 => '1',
            I5 => '1',
            O5 => gen(t+1)(word_size-3-t+1), 
            O6 => prop(t+1)(word_size-3-t+1)
            );
        B_steps1: 
            for w in 0 to ((word_size-3)-(word_size-3-t)) generate
                lut_inst4: lut6_2 
                    generic map(INIT => X"7878787880808080")
                    port map(
                    I0 => b(t + 2),
                    I1 => a((word_size-3-t)+2+w), --------------- $$$$$$$$$$$$$$ changes have been made here. instead of + 1 i have changed it to + 2
                    I2 => chain(t)(word_size-3-t+2+1+w),
                    I3 => '1',
                    I4 => '1',
                    I5 => '1',
                    O5 => gen(t+1)(word_size-3-t+1+1+w), 
                    O6 => prop(t+1)(word_size-3-t+1+1+w)
                    );
        end generate B_steps1;
        ----------------------------------------------------------
        carry_chain_type_B:										
            for x in 0 to (word_size/4 - 1) generate			
                carry_inst1: CARRY4										
                port map (													
                    DI => gen(t+1)(x*4+3 downto x*4), 											
                    S => prop(t+1)(x*4+3 downto x*4),											
                    O => output(t+1)(x*4+3 downto x*4),										
                    CO => carries(t+1)(x*4+3 downto x*4),										
                    CI => input_carry(t+1)(x),							
                    CYINIT => '0'
                    );
                input_carry(t+1)(x+1) <= carries(t+1)(x*4+3);
                
                chain(t+1)(x*4+3 downto x*4) <= output(t+1)(x*4+3 downto x*4);
        end generate carry_chain_type_B;
        chain(t+1)(word_size) <= carries(t+1)(word_size-1);
    end generate Type_B;
    -----------------------------------------
    -- Another LUT Type B
    lut_inst5: lut6_2 
        generic map(INIT => X"7878787880808080")
        port map(
        I0 => b(word_size-1),
        I1 => a(0),
        I2 => chain(word_size-3)(1),
        I3 => '1',
        I4 => '1',
        I5 => '1',
        O5 => gen(word_size-2)(0), 
        O6 => prop(word_size-2)(0)
        );
    ------------------------------------------
    -- A single LUT type C
    lut_inst6: lut6_2 
    generic map(INIT => X"F0F0F0F088888888")
    port map(
    I0 => b(0),
    I1 => a(0),
    I2 => chain(word_size-3)(2),
    I3 => '1',
    I4 => '1',
    I5 => '1',
    O5 => prod(0), ------------------- product bit 0
    O6 => prop(word_size-2)(1)
    );
    gen(word_size-2)(1) <= '0'; -- $$$$$$$$$$$$$$$$$$
    -----------------------------------------------
    last_row: 
    for y in 0 to word_size - 3 generate
        lut_inst7: lut6_2 
            generic map(INIT => X"7878787880808080")
            port map(
            I0 => b(word_size-1),
            I1 => a(y+2),
            I2 => chain(word_size-3)(3+y),
            I3 => '1',
            I4 => '1',
            I5 => '1',
            O5 => gen(word_size-2)(y+2), 
            O6 => prop(word_size-2)(y+2)
            );
    end generate last_row;
    last_row_carry_chain:
    for z in 0 to (word_size/4 - 1) generate			
        carry_inst2: CARRY4										
        port map (													
            DI => gen(word_size-2)(z*4+3 downto z*4), 											
            S => prop(word_size-2)(z*4+3 downto z*4),											
            O => output(word_size-2)(z*4+3 downto z*4),								
            CO => carries(word_size-2)(z*4+3 downto z*4),										
            CI => input_carry(word_size-2)(z),							
            CYINIT => '0'
            );
        input_carry(word_size-2)(z+1) <= carries(word_size-2)(z*4+3);
        chain(word_size-2)(z*4+3 downto z*4) <= output(word_size-2)(z*4+3 downto z*4);
    end generate last_row_carry_chain;
    chain(word_size-2)(word_size) <= carries(word_size-2)(word_size-1);
    
    assign_values:
    for zz in 0 to (word_size - 2) generate
        prod(zz+1) <= chain(zz)(0);
    end generate assign_values;
    prod((word_size*2)-1 downto word_size) <= chain(word_size-2)(word_size downto 1);
   
end Behavioral;