

��MySQL���û����ݷ���mysql���user����ڸñ��б������û����û����������Ȩ����Ϣ��

## �����û�

���`CREATE USER 'username'@'host' IDENTIFIED BY 'password'`;

˵����
- username����¼���û���
- password���ǵ�¼������
- host��ָ�����Ե�¼������������localhost��ʾ������%��ʾ��������

������

    CREATE USER 'testuser'@'%' IDENTIFIED BY '123';
    
����һ���û���user���л����һ���µ����ݣ����Ǹ��û���û���κ�Ȩ�޵�

## �����û�Ȩ��
**����û�Ȩ��**

���`GRANT privileges ON databasename.tablename TO 'username'@'host'`
˵����

- privileges���û��Ĳ���Ȩ��,��SELECT , INSERT , UPDATE�����Ҫ��������Ȩ����all privileges
- databasename�����ݿ�����ָ��Ҫ����Ȩ�޵Ŀ⣬���Ҫ�����п�����Ȩ����`*`����
- tablename��������ָ��Ҫ����Ȩ�޵ı����Ҫ�����б�����Ȩ����`*`����
- username��������Ȩ�޵��û���
- host��������Ȩ�޵����������Ҫ��������������Ȩ����%����

������

    # ���û�testuser���������������϶�test���SELECTȨ��
    GRANT SELECT ON test.* TO 'testuser'@'%';
    
    # �����пⶼ����INSERTȨ�޺�user���еĸ��û�Insert_priv��ֵ����ΪY
    GRANT INSERT ON *.* TO 'testuser'@'%';
    
    # �����û���������Ȩ��
    GRANT ALL privileges ON *.* TO 'testuser'@'%';
    
    # ִ��������Ȩ�޵�����󣬱���ִ����������ʹ�޸���Ч
    flush privileges;

**�����û�Ȩ��**

���`REVOKE privileges ON databasename.tablename FROM 'username'@'host';`

���ӣ�`REVOKE ALL privileges ON *.* FROM 'testuser'@'%';`

����Ȩ�޵�ʹ�÷���������Ȩ�����ƣ���������ķ�ʽ���д����ɡ�

## �޸��û�����

����: `SET PASSWORD FOR 'username'@'host' = PASSWORD('newpassword');`

����: `SET PASSWORD FOR 'testuser'@'%' = PASSWORD('abcdef');`

## �鿴Ȩ��

�鿴���û���Ȩ�ޣ�`SHOW GRANTS;`

�鿴ָ���û���Ȩ�ޣ�`SHOW GRANTS FOR 'username'@'host';`

������

    # ��ѯ�û�testuser�����������ϵ�Ȩ��
    SHOW GRANTS FOR 'testuser'@'%'

## ɾ���û�

ɾ���û���`DROP USER 'username'@'host';`

������`DROP USER 'testuser'@'%';`

�û���ɾ����user���оͲ����ڸ��û������ݡ�������һ����Ҫע����ǣ�user������ͨ���û�������������Ψһȷ��һ���û����������`CREATE USER 'testuser'@'10.1.1.1' IDENTIFIED BY 'password';`��`CREATE USER 'testuser'@'10.1.1.2' IDENTIFIED BY 'password';`��������䣬����user���в���������Ϣ����ˣ�ɾ���û���ʱ��Ҳֻ�ܷ�����ɾ����
    

## �ܽ�

��ƪ���½�����MySQL�����û����õ�һЩ������������û�������û�Ȩ�ޣ�ɾ���û�Ȩ�ޣ��޸��û����룬��ѯ�û�Ȩ�ޡ����⣬user��Ҳ��һ�����ݱ����Ҳ������MySQL�������SQL���ȥ�����ñ�����ɾ���û�����ʹ��`DELETE FORM USER WHERE User = testuser`����Ȼ�����ַ�ʽ�ǲ����Ƽ��ģ����ڸ������������ա�

ʱ�䣺2019��01��08��