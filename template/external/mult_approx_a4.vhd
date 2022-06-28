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

type data_type is array(11 downto 0) of std_logic_vector(11 downto 0);
signal prop: data_type;
signal gen: data_type;
signal carries: data_type;
signal output: data_type;
type input_carries is array((word_size - 2 ) downto 0) of std_logic_vector(word_size/4 downto 0);
signal input_carry: data_type; --input_carries;
type interm_results is array((word_size - 2) downto 0) of std_logic_vector(word_size downto 0);
signal chain: data_type; --interm_results;
signal anding1 : std_logic;
signal anding2 : std_logic;
signal anding3 : std_logic;
--signal temp : std_logic;
signal temp1, temp2, temp3, temp4, temp : std_logic;
signal prod1: STD_LOGIC_VECTOR (word_size + 3 downto 0);
signal prod2: STD_LOGIC_VECTOR (word_size + 3 downto 0);

begin

 set_initial_carry:
    for car in 0 to (word_size - 2) generate
        input_carry(car)(0) <= '0';
    end generate set_initial_carry; 
    
  Type_1: 
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
           end generate Type_1;
           
           lut_inst1: lut6_2 
               generic map(INIT => X"88888888F000F000")
               port map(
               I0 => b(0), 
               I1 => a(0), 
               I2 => b(1), 
               I3 => a(word_size-1), 
               I4 => '1',
               I5 => '1',
               O5 => temp1 ,
               O6 => prod(0)
               );
               
               gen(0)(3) <= '0';
               prop(0)(3) <= temp1;	
               
carry_chain_1:
               for z in 0 to 0 generate            
               carry_inst2: CARRY4                                        
               port map (                                                    
               DI => gen(0)(z*4+3 downto z*4),                                             
               S => prop(0)(z*4+3 downto z*4),                                            
               O => output(0)(z*4+3 downto z*4),                                
               CO => carries(0)(z*4+3 downto z*4),                                        
               CI => input_carry(0)(z),                            
               CYINIT => '0'
               );
               input_carry(0)(z+1) <= carries(0)(z*4+3);
               chain(0)(z*4+3 downto z*4) <= output(0)(z*4+3 downto z*4);
               end generate carry_chain_1;
               chain(0)(4) <= carries(0)(3);     
               
     Type_2: 
                              for i in 0 to (word_size - 2) generate 
                              lut_inst0: lut6_2 
                              generic map(INIT => X"7888788880008000")
                              port map(
                              I0 => b(2), --h
                              I1 => a(i+1), --c
                              I2 => b(3), -- g
                              I3 => a(i), -- d
                              I4 => '1',
                              I5 => '1',
                              O5 => gen(1)(i),
                              O6 => prop(1)(i)
                              );
                              end generate Type_2;           
                                   
                         lut_inst2: lut6_2 
                              generic map(INIT => X"88888888F000F000")
                              port map(
                              I0 => b(2), 
                              I1 => a(0), 
                              I2 => b(3), 
                              I3 => a(word_size-1), 
                              I4 => '1',
                              I5 => '1',
                              O5 => temp2 ,
                              O6 => prod2(0)
                              );
                                                     
                              gen(1)(3) <= '0';
                              prop(1)(3) <= temp2;               
                              
     carry_chain_2:
                              for z in 0 to 0 generate            
                              carry_inst2: CARRY4                                        
                              port map (                                                    
                              DI => gen(1)(z*4+3 downto z*4),                                             
                              S => prop(1)(z*4+3 downto z*4),                                            
                              O => output(1)(z*4+3 downto z*4),                                
                              CO => carries(1)(z*4+3 downto z*4),                                        
                              CI => input_carry(1)(z),                            
                              CYINIT => '0'
                              );
                              input_carry(1)(z+1) <= carries(1)(z*4+3);
                              chain(1)(z*4+3 downto z*4) <= output(1)(z*4+3 downto z*4);
                              end generate carry_chain_2;
                              chain(1)(4) <= carries(1)(3);                              

add: for i in 0 to 1 generate
               lut_inst3: lut6_2 
                         generic map(INIT => X"9666666666666666")
                         port map(
                         I0 => chain(1)(i), 
                         I1 => chain(0)(i+2), 
                         I2 => b(1), 
                         I3 => a(i+1), 
                         I4 => b(2),
                         I5 => a(i),
                         O5 => temp,
                         O6 => prod(i+3)
                         );               
               end generate add;
               
lut_inst4: lut6_2 
                             generic map(INIT => X"666666660FF00FF0")
                             port map(
                             I0 => prod2(0), 
                             I1 => chain(0)(1), 
                             I2 => chain(1)(2), 
                             I3 => chain(0)(4), 
                             I4 => '1',
                             I5 => '1',
                             O5 => prod(5) ,
                             O6 => prod(2)
                            );  
                         prod(1) <= chain(0)(0);
                         prod(6) <= chain(1)(3);
                         prod(7) <= chain(1)(4);               

end Behavioral;
